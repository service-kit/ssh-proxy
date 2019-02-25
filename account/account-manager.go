package account

import "sync"

var mgr *AccountManager
var once sync.Once

type AccountManager struct {
	accountMap map[string]string
}

func GetInstance() *AccountManager {
	once.Do(func() {
		mgr = new(AccountManager)
	})
	return mgr
}

func (a *AccountManager) InitManager() error {
	a.accountMap = make(map[string]string)
	return nil
}

func (a *AccountManager) RegisterAccount(user, password string) bool {
	if "" == user || "" == password {
		return false
	}
	tp := a.accountMap[user]
	if "" == tp {
		a.accountMap[user] = password
		return true
	}
	return false
}

func (a *AccountManager) CheckUserPassword(user, password string) bool {
	if "" == user || "" == password {
		return false
	}
	tp := a.accountMap[user]
	if "" == tp {
		return false
	}
	return password == tp
}
