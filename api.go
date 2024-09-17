package wildlifenl

import (
	"database/sql"
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
	"github.com/UtrechtUniversity/wildlifenl/timeseries"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/patrickmn/go-cache"
)

const appName = "WildlifeNL"
const appDescription = "This is the WildlifeNL API. Before you can start making calls to the provided end-points you should acquire a bearer token. To do so, make a POST request at /auth/ providing the required fields including a valid email address. A validation code will be send to that email address. Then, make a PUT request at /auth/ providing the same email address and the validation code. The response will include a field named \"token\" containing your bearer token. Use this bearer token in the header of any future calls you make."

var (
	configuration *Configuration
	relationalDB  *sql.DB
	timeseriesDB  *timeseries.Timeseries
	sessions      *cache.Cache
	authRequests  *cache.Cache
)

func Start(config *Configuration) error {
	configuration = config
	if err := loadRelationalDB(configuration); err != nil {
		return fmt.Errorf("could not connect to relational database: %w", err)
	}
	if err := loadTimeseriesDB(configuration); err != nil {
		return fmt.Errorf("could not connect to timeseries database: %w", err)
	}
	sessions = cache.New(time.Duration(configuration.CacheSessionDurationMinutes)*time.Minute, 12*time.Hour)
	authRequests = cache.New(time.Duration(configuration.CacheAuthRequestDurationMinutes)*time.Minute, 12*time.Hour)
	apiConfig := huma.DefaultConfig(appName, config.Version)
	apiConfig.Info.Description = appDescription
	apiConfig.Security = []map[string][]string{{"auth": {}}}
	apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{"auth": {Type: "http", Scheme: "bearer"}}
	apiConfig.DocsPath = "/"

	router := http.NewServeMux()
	api := humago.New(router, apiConfig)
	api.UseMiddleware(NewAuthMiddleware(api))
	huma.AutoRegister(api, newAnimalOperations(relationalDB))
	huma.AutoRegister(api, newAnswerOperations(relationalDB))
	huma.AutoRegister(api, newAuthOperations(relationalDB))
	huma.AutoRegister(api, newBorneSensorDeploymentOperations())
	huma.AutoRegister(api, newBorneSensorReadingOperations())
	huma.AutoRegister(api, newExperimentOperations(relationalDB))
	huma.AutoRegister(api, newInteractionOperations(relationalDB))
	huma.AutoRegister(api, newInteractionTypeOperations(relationalDB))
	huma.AutoRegister(api, newLivingLabOperations(relationalDB))
	huma.AutoRegister(api, newMeOperations(relationalDB))
	huma.AutoRegister(api, newQuestionOperations(relationalDB))
	huma.AutoRegister(api, newQuestionnaireOperations(relationalDB))
	huma.AutoRegister(api, newRoleOperations(relationalDB))
	huma.AutoRegister(api, newSpeciesOperations(relationalDB))
	huma.AutoRegister(api, newUserOperations(relationalDB))
	huma.AutoRegister(api, newZoneOperations(relationalDB))
	return http.ListenAndServe(configuration.Host+":"+strconv.Itoa(configuration.Port), router)
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
	timeseriesDB = timeseries.NewTimeseries(config.TimeseriesDatabaseURL, config.TimeseriesDatabaseOrganization, config.TimeseriesDatabaseToken)
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
	if err := sendCodeByEmail(displayNameApp, displayNameUser, email, code); err != nil {
		return err
	}
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
		return nil, nil
	}
	authenticationRequest := c.(AuthenticationRequest)
	if authenticationRequest.code != code {
		return nil, nil
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
		log.Fatal("ERROR getting credential from store:", err)
	}
	if credential == nil {
		return nil
	}
	sessions.SetDefault(token, credential)
	return credential
}
