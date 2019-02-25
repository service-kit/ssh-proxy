package common

const (
	CODE    = "code"
	ERROR   = "error"
	TOKEN   = "token"
	SUCCESS = "success"
	FAIL    = "fail"
)

const (
	ERROR_VERIFY_NOT_PASS  = "verify not pass"
	ERROR_CAN_NOT_REGISTER = "can not register"
	ERROR_EXIST            = "register exist"
)

const (
	SWITHC_ON  = 1
	SWITHC_OFF = 0
)

const (
	LL_DEBUG = "debug"
	LL_INFO  = "info"
	LL_ERROR = "error"
)

type SSHRegisterRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type SSHRegisterResponse struct {
	Result string `json:"result"`
	Code   string `json:"code"`
}

type SSHLoginRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type SSHLoginResponse struct {
	Result string `json:"result"`
	Code   string `json:"code"`
	Token  string `json:"token"`
}

type SSHCmdRequest struct {
	UserId   string `json:"user_id"`
	Host     string `json:"host"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Cmd      string `json:"cmd"`
}

const (
	RSP_CODE_PARAM_INVALID    = "param invalid"
	RSP_CODE_NOT_ACCOUNT      = "not account"
	RSP_CODE_PASSWORD_ERR     = "password err"
	RSP_CODE_DO_CMD_ERR       = "do cmd err"
	RSP_CODE_CREATE_PROXY_ERR = "create proxy err"
)

type SSHCmdResponse struct {
	Out    string `json:"out"`
	Result string `json:"result"`
	Code   string `json:"code"`
}
