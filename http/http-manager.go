package http

import (
	"github.com/service-kit/ssh-proxy/config"
	"github.com/service-kit/ssh-proxy/log"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type HttpManager struct {
	addr string
	wg   *sync.WaitGroup
}

var m *HttpManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *HttpManager {
	once.Do(func() {
		m = &HttpManager{}
	})
	return m
}

func (self *HttpManager) InitManager(wg *sync.WaitGroup) error {
	logger = log.GetInstance().GetLogger()
	self.wg = wg
	var err error = nil
	self.addr, err = config.GetInstance().GetConfig("SSH_PROXY_HTTP_ADDR")
	if nil != err {
		return err
	}
	self.wg.Add(1)
	http.HandleFunc("/register", handleSSHRegisterRequest)
	http.HandleFunc("/login", handleSSHLoginRequest)
	http.HandleFunc("/cmd", handleSSHCmdRequest)
	go self.startHttpServer()
	return nil
}

func (self *HttpManager) startHttpServer() {
	logger.Info("Start Http Server", zap.String("addr", self.addr))
	defer self.wg.Done()
	http.ListenAndServe(self.addr, nil)
}
