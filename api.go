package wildlifenl

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/patrickmn/go-cache"
)

const addr = "localhost:8888" // TODO configuration

var app *App

func Start() error {
	db, err := newStore()
	if err != nil {
		return fmt.Errorf("could not create database connection: %w", err)
	}
	app = &App{
		sessions: cache.New(4*time.Hour, 12*time.Hour),
		store:    db,
	}
	config := huma.DefaultConfig("WildlifeNL", "0.1")
	config.Security = []map[string][]string{{"auth": {}}}
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"auth": {Type: "http", Scheme: "bearer"},
	}
	router := http.NewServeMux()
	api := humago.New(router, config)
	api.UseMiddleware(NewAuthMiddleware(api))
	router.HandleFunc("/", rootHandlerFunc)
	huma.AutoRegister(api, &animalOperations{})
	return http.ListenAndServe(addr, router)
}

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		account := app.CreateAccount()
		if account != nil {
			w.Write([]byte(account.Token))
		}
	}
}

func NewAuthMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		var anyOfNeededScopes []string
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var ok bool
			if anyOfNeededScopes, ok = opScheme["auth"]; ok {
				isAuthorizationRequired = true
				break
			}
		}
		if !isAuthorizationRequired {
			next(ctx)
			return
		}
		token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
		if len(token) == 0 {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing bearer token.")
			return
		}
		account := app.GetAccount(token)
		if account == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid bearer token.")
			return
		}
		for _, scope := range account.Scopes {
			if slices.Contains(anyOfNeededScopes, scope) {
				next(ctx)
				return
			}
		}
		huma.WriteErr(api, ctx, http.StatusForbidden, "Bearer token is not authorized for any of the scopes: "+strings.Join(anyOfNeededScopes, ", "))
	}
}
