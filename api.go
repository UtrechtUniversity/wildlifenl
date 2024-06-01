package wildlifenl

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/patrickmn/go-cache"
	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"

	_ "github.com/lib/pq" // postgresql
)

const appName = "WildlifeNL"
const appVersion = "0.1"

var (
	database     *sql.DB
	sessions     *cache.Cache
	authRequests *cache.Cache
)

type Configuration struct {
	Host      string
	Port      int
	DbHost    string
	DbName    string
	DbUser    string
	DbPass    string
	DbSSLmode string
}

type AuthenticationRequest struct {
	appName  string
	userName string
	email    string
	code     string
}

type Input struct {
	credential *models.Credential
}

func (m *Input) Resolve(ctx huma.Context) []error {
	token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
	m.credential = getCredential(token)
	return nil
}

func Start(config *Configuration) error {
	connStr := "postgres://" + config.DbUser + ":" + config.DbPass + "@" + config.DbHost + "/" + config.DbName + "?sslmode=" + config.DbSSLmode
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	db.SetConnMaxLifetime(1 * time.Hour)
	database = db
	sessions = cache.New(4*time.Hour, 12*time.Hour)
	authRequests = cache.New(1*time.Hour, 12*time.Hour)
	apiConfig := huma.DefaultConfig(appName, appVersion)
	apiConfig.Security = []map[string][]string{{"auth": {}}}
	apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{"auth": {Type: "http", Scheme: "bearer"}}
	router := http.NewServeMux()
	api := humago.New(router, apiConfig)
	api.UseMiddleware(NewAuthMiddleware(api))
	huma.AutoRegister(api, new(authOperations))
	huma.AutoRegister(api, new(animalOperations))
	huma.AutoRegister(api, new(meOperations))
	huma.AutoRegister(api, new(noticeOperations))
	huma.AutoRegister(api, new(noticeTypeOperations))
	huma.AutoRegister(api, new(speciesOperations))
	huma.AutoRegister(api, new(userOperations))
	return http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), router)
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
		account := getCredential(token)
		if account == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid bearer token.")
			return
		}
		if len(anyOfNeededScopes) == 0 {
			next(ctx)
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

func authenticate(displayNameApp, displayNameUser, email string) error {
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.IntN(10))
	}

	// TODO send email message here
	// sendEmail(request)

	authenticationRequest := AuthenticationRequest{
		appName:  displayNameApp,
		userName: displayNameUser,
		email:    email,
		code:     code,
	}
	authRequests.SetDefault(email, authenticationRequest)
	return nil
}

func authorize(email, code string) (*models.Credential, error) {
	c, ok := authRequests.Get(email)
	if !ok {
		return nil, errors.New("email address not found")
	}
	authenticationRequest := c.(AuthenticationRequest)
	if authenticationRequest.code != code {
		return nil, errors.New("code does not match")
	}
	account := stores.NewCredentialStore(database).Create(authenticationRequest.appName, authenticationRequest.userName, authenticationRequest.email)
	sessions.SetDefault(account.Token, account)
	return account, nil
}

func getCredential(token string) *models.Credential {
	if credential, ok := sessions.Get(token); ok {
		return credential.(*models.Credential)
	}
	credential := stores.NewCredentialStore(database).Get(token)
	if credential == nil {
		return nil
	}
	sessions.SetDefault(token, credential)
	return credential
}

func scopesAsMarkdown(input []string) string {
	result := make([]string, 0)
	for _, value := range input {
		result = append(result, "`"+value+"`")
	}
	return strings.Join(result, ", ")
}

func logTrace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File + " " + strconv.Itoa(frame.Line) + " " + frame.Function + ": "
}
