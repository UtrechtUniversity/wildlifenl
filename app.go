package wildlifenl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

const addr = "localhost:80" // TODO configuration

var app *App

type App struct {
	sessions *cache.Cache
	database *Database
}

func Start() error {
	db, err := newDataBase()
	if err != nil {
		return fmt.Errorf("could not create database connection: %w", err)
	}
	app = &App{
		sessions: cache.New(4*time.Hour, 12*time.Hour),
		database: db,
	}
	http.HandleFunc("/", app.rootHandler)
	http.HandleFunc("/me/", newMe().handle)
	return http.ListenAndServe(addr, nil)
}

func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var email string
		var password string
		var ok bool
		if email, password, ok = r.BasicAuth(); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		account := a.database.getAccount(email, password)
		if account == nil {
			account = a.database.addAccount(email, password)
			if account == nil {
				w.WriteHeader(http.StatusUnauthorized) // maybe there is a better status for "account exists but password is wrong"?
				return
			}
		}
		token := uuid.New().String()
		a.sessions.Set(token, account, cache.DefaultExpiration)
		result := make(map[string]string)
		result["bearer"] = token
		data, _ := json.Marshal(result)
		writeResponseJSON(w, data)
	}
}

func (a *App) authenticate(r *http.Request) bool {
	return a.getAuthorization(r) != nil
}

func (a *App) getAuthorization(r *http.Request) *Account {
	token := getBearerToken(r)
	if value, ok := a.sessions.Get(token); ok {
		return value.(*Account)
	}
	return nil
}
