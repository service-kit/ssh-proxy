package http

import (
	"encoding/json"
	"github.com/service-kit/ssh-proxy/account"
	"github.com/service-kit/ssh-proxy/common"
	"github.com/service-kit/ssh-proxy/ssh"
	"go.uber.org/zap"
	"net/http"
)

func handleSSHRegisterRequest(w http.ResponseWriter, r *http.Request) {
	registerReq := new(common.SSHRegisterRequest)
	json.NewDecoder(r.Body).Decode(registerReq)
	dat, _ := json.Marshal(registerReq)
	logger.Info(string(dat))
	res := account.GetInstance().RegisterAccount(registerReq.UserId, registerReq.Password)
	rsp := new(common.SSHRegisterResponse)
	if res {
		rsp.Result = common.SUCCESS
	} else {
		rsp.Result = common.FAIL
	}
	json.NewEncoder(w).Encode(rsp)
}

func handleSSHLoginRequest(w http.ResponseWriter, r *http.Request) {
	loginReq := new(common.SSHLoginRequest)
	json.NewDecoder(r.Body).Decode(loginReq)
	res := account.GetInstance().CheckUserPassword(loginReq.UserId, loginReq.Password)
	rsp := new(common.SSHRegisterResponse)
	if res {
		rsp.Result = common.SUCCESS
	} else {
		rsp.Result = common.FAIL
	}
	json.NewEncoder(w).Encode(rsp)
}

func handleSSHCmdRequest(w http.ResponseWriter, r *http.Request) {
	cmdReq := new(common.SSHCmdRequest)
	rsp := new(common.SSHCmdResponse)
	defer func() {
		json.NewEncoder(w).Encode(rsp)
	}()
	rsp.Result = common.FAIL
	err := json.NewDecoder(r.Body).Decode(cmdReq)
	if nil != err {
		logger.Error("do shell cmd err", zap.Error(err))
		rsp.Code = common.RSP_CODE_PARAM_INVALID
		return
	}
	p := ssh.GetInstance().GetProxy(cmdReq.UserId)
	if nil == p {
		if "" == cmdReq.Host || "" == cmdReq.UserName {
			logger.Error("param invalid", zap.String("host", cmdReq.Host), zap.String("user_name", cmdReq.UserName))
			rsp.Code = common.RSP_CODE_PARAM_INVALID
			return
		}
		p, err = ssh.GetInstance().CreateSSHProxy(cmdReq.UserId, cmdReq.Host, cmdReq.UserName, cmdReq.Password, cmdReq.Port)
		if nil != err {
			logger.Error("do shell cmd err", zap.Error(err))
			rsp.Code = common.RSP_CODE_CREATE_PROXY_ERR
			return
		}
	}
	logger.Info("recv cmd req", zap.String("user id", cmdReq.UserId), zap.String("host", p.Host()), zap.String("cmd", cmdReq.Cmd))
	res, err := p.DoCmd(cmdReq.Cmd)
	if nil != err {
		logger.Error("do shell cmd err", zap.Error(err))
		rsp.Code = common.RSP_CODE_DO_CMD_ERR
		return
	}
	rsp.Result = common.SUCCESS
	rsp.Out = res
}
