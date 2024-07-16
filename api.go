package wildlifenl

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/patrickmn/go-cache"
)

const appName = "WildlifeNL"
const appVersion = "0.1"

var (
	relationalDB *sql.DB
	timeseriesDB *stores.Timeseries
	sessions     *cache.Cache
	authRequests *cache.Cache
)

func Start(config *Configuration) error {
	if err := loadRelationalDB(config); err != nil {
		return fmt.Errorf("could not connect to relational database: %w", err)
	}
	if err := loadTimeseriesDB(config); err != nil {
		return fmt.Errorf("could not connect to timeseries database: %w", err)
	}
	sessions = cache.New(time.Duration(config.CacheSessionDurationMinutes)*time.Minute, 12*time.Hour)
	authRequests = cache.New(time.Duration(config.CacheAuthRequestDurationMinutes)*time.Minute, 12*time.Hour)
	apiConfig := huma.DefaultConfig(appName, appVersion)
	apiConfig.Security = []map[string][]string{{"auth": {}}}
	apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{"auth": {Type: "http", Scheme: "bearer"}}
	apiConfig.DocsPath = "/"
	router := http.NewServeMux()
	api := humago.New(router, apiConfig)
	api.UseMiddleware(NewAuthMiddleware(api))
	huma.AutoRegister(api, newAnimalOperations(relationalDB))
	huma.AutoRegister(api, newAreaOperations(relationalDB))
	huma.AutoRegister(api, newAuthOperations(relationalDB))
	huma.AutoRegister(api, newBorneSensorDeploymentOperations())
	huma.AutoRegister(api, newBorneSensorReadingOperations())
	huma.AutoRegister(api, newInteractionOperations(relationalDB))
	huma.AutoRegister(api, newMeOperations(relationalDB))
	huma.AutoRegister(api, newNoticeOperations(relationalDB))
	huma.AutoRegister(api, newNoticeTypeOperations(relationalDB))
	huma.AutoRegister(api, newLivingLabOperations(relationalDB))
	huma.AutoRegister(api, newRoleOperations(relationalDB))
	huma.AutoRegister(api, newSpeciesOperations(relationalDB))
	huma.AutoRegister(api, newTrackingEventOperations(relationalDB))
	huma.AutoRegister(api, newUserOperations(relationalDB))
	return http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), router)
}

func loadRelationalDB(config *Configuration) error {
	connStr := "postgres://" + config.RelationalDatabaseUser + ":" + config.RelationalDatabasePass + "@" + config.RelationalDatabaseHost
	if config.RelationalDatabasePort > 0 {
		connStr += ":" + strconv.Itoa(config.RelationalDatabasePort)
	}
	connStr += "/" + config.RelationalDatabaseName + "?sslmode=" + config.RelationalDatabaseSSLmode
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(1 * time.Hour)
	relationalDB = db
	return relationalDB.Ping()
}

func loadTimeseriesDB(config *Configuration) error {
	timeseriesDB = stores.NewTimeseries(config.TimeseriesDatabaseURL, config.TimeseriesDatabaseOrganization, config.TimeseriesDatabaseToken)
	return timeseriesDB.Ping()
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
		return nil, errors.New("email address not found, possibly the code has expired")
	}
	authenticationRequest := c.(AuthenticationRequest)
	if authenticationRequest.code != code {
		return nil, errors.New("provided code does not match the sent code")
	}
	account, err := stores.NewCredentialStore(relationalDB).Create(authenticationRequest.appName, authenticationRequest.userName, authenticationRequest.email)
	if err != nil {
		return nil, err
	}
	sessions.SetDefault(account.Token, account)
	return account, nil
}

func getCredential(token string) *models.Credential {
	if credential, ok := sessions.Get(token); ok {
		return credential.(*models.Credential)
	}
	credential, err := stores.NewCredentialStore(relationalDB).Get(token)
	if err != nil {
		log.Fatal("ERROR when using credential store:", err)
	}
	if credential == nil {
		return nil
	}
	sessions.SetDefault(token, credential)
	return credential
}
