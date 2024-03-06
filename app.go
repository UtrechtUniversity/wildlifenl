package wildlifenl

import (
	"github.com/patrickmn/go-cache"
)

type Account struct {
	Token  string
	Scopes []string
}

type App struct {
	sessions *cache.Cache
	store    *Store
}

func (a *App) CreateAccount() *Account {
	account := a.store.createAccount()
	a.sessions.SetDefault(account.Token, account)
	return account
}

func (a *App) GetAccount(token string) *Account {
	if account, ok := a.sessions.Get(token); ok {
		return account.(*Account)
	}
	account := a.store.getAccount(token)
	if account == nil {
		return nil
	}
	a.sessions.SetDefault(token, account)
	return account
}
