package ssh

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"time"
)

type SSHProxy struct {
	session   *ssh.Session
	userName  string
	password  string
	host      string
	port      int
	stdoutBuf io.Reader
	stdinBuf  io.WriteCloser
}

func (sp *SSHProxy) Init(user, password, host, key string, port int) (err error) {
	sp.userName = user
	sp.password = password
	sp.host = host
	sp.port = port
	return
}

func (sp *SSHProxy) connect(user, password, host, key string, port int, cipherList []string) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(cipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "seaofstars.developer@gmail.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: cipherList,
		}
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}

func (sp *SSHProxy) DoCmd(cmd string) (string, error) {
	var err error
	sp.session, err = sp.connect(sp.userName, sp.password, sp.host, "", sp.port, nil)
	if nil != err {
		return "", err
	}
	defer sp.session.Close()
	sp.stdoutBuf, _ = sp.session.StdoutPipe()
	sp.stdinBuf, _ = sp.session.StdinPipe()
	if nil == sp.session {
		return "", errors.New("ssh not connected")
	}
	sp.session.Run(cmd)
	bytes, err := ioutil.ReadAll(sp.stdoutBuf)
	if nil != err {
		return "", err
	}
	res := string(bytes)
	return res, nil
}

func (sp *SSHProxy) Host() string {
	return sp.host
}
