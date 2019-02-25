package ssh

import (
	"sync"
)

var mgr *SSHManager
var once sync.Once

type SSHManager struct {
	sshProxyMap map[string]*SSHProxy
}

func GetInstance() *SSHManager {
	once.Do(func() {
		mgr = new(SSHManager)
	})
	return mgr
}

func (s *SSHManager) InitManager() error {
	s.sshProxyMap = make(map[string]*SSHProxy)
	return nil
}

func (s *SSHManager) CreateSSHProxy(userId, host, user, passwd string, port int) (*SSHProxy, error) {
	sp := new(SSHProxy)
	err := sp.Init(user, passwd, host, "", port)
	if nil != err {
		return nil, err
	}
	s.sshProxyMap[userId] = sp
	return sp, nil
}

func (s *SSHManager) GetProxy(id string) *SSHProxy {
	return s.sshProxyMap[id]
}
