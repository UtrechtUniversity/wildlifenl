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
	"github.com/UtrechtUniversity/wildlifenl/timeseries"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

var (
	relationalDB *sql.DB
	timeseriesDB *timeseries.Timeseries
	mailer       *Mailer
	sessions     *cache.Cache
	authRequests *cache.Cache
)

func Start(config *Configuration) error {
	if err := initializeRelationalDB(config); err != nil {
		return fmt.Errorf("could not initialize Relational database: %w", err)
	}
	if err := initializeTimeseriesDB(config); err != nil {
		return fmt.Errorf("could not initiliaze Timeseries database: %w", err)
	}
	if err := initializeMailer(config); err != nil {
		return fmt.Errorf("could not initialize Mailer service: %w", err)
	}
	if err := timeseriesDB.CreateBucketIfNotExists("animals"); err != nil {
		return fmt.Errorf("could not create Timeseries bucket: %w", err)
	}
	if err := timeseriesDB.CreateBucketIfNotExists("humans"); err != nil {
		return fmt.Errorf("could not create Timeseries bucket: %w", err)
	}
	if err := assertAdminUserExists(config); err != nil {
		return fmt.Errorf("could not asset that admin user exists: %w", err)
	}

	sessions = cache.New(time.Duration(config.CacheSessionDurationMinutes)*time.Minute, 12*time.Hour)
	authRequests = cache.New(time.Duration(config.CacheAuthRequestDurationMinutes)*time.Minute, 12*time.Hour)
	apiConfig := huma.DefaultConfig(appName, config.Version)
	apiConfig.Info.Description = appDescription
	apiConfig.Security = []map[string][]string{{"auth": {}}}
	apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{"auth": {Type: "http", Scheme: "bearer"}}
	apiConfig.DocsPath = "/"

	router := NewServeMux()
	api := humago.New(router, apiConfig)
	api.UseMiddleware(NewMiddleware(api))
	huma.AutoRegister(api, newAlarmOperations())
	huma.AutoRegister(api, newAnimalOperations())
	huma.AutoRegister(api, newAnswerOperations())
	huma.AutoRegister(api, newAuthOperations())
	huma.AutoRegister(api, newBorneSensorDeploymentOperations())
	huma.AutoRegister(api, newBorneSensorReadingOperations())
	huma.AutoRegister(api, newConveyanceOperations())
	huma.AutoRegister(api, newDetectionOperations())
	huma.AutoRegister(api, newExperimentOperations())
	huma.AutoRegister(api, newInteractionOperations())
	huma.AutoRegister(api, newInteractionTypeOperations())
	huma.AutoRegister(api, newLivingLabOperations())
	huma.AutoRegister(api, newMessageOperations())
	huma.AutoRegister(api, newProfileOperations())
	huma.AutoRegister(api, newQuestionOperations())
	huma.AutoRegister(api, newQuestionnaireOperations())
	huma.AutoRegister(api, newResponseOperations())
	huma.AutoRegister(api, newRoleOperations())
	huma.AutoRegister(api, newSpeciesOperations())
	huma.AutoRegister(api, newTrackingReadingOperations())
	huma.AutoRegister(api, newUserOperations())
	huma.AutoRegister(api, newZoneOperations())
	return http.ListenAndServe(config.Host+":"+strconv.Itoa(config.Port), router)
}

func initializeRelationalDB(config *Configuration) error {
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

func initializeTimeseriesDB(config *Configuration) error {
	timeseriesDB = timeseries.NewTimeseries(config.TimeseriesDatabaseURL, config.TimeseriesDatabaseOrganization, config.TimeseriesDatabaseToken)
	return timeseriesDB.Ping()
}

func initializeMailer(config *Configuration) error {
	mailer = newMailer(config)
	return mailer.Ping()
}

func NewMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
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
		credential := getCredential(token)
		if credential == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid bearer token.")
			return
		}
		if len(anyOfNeededScopes) == 0 {
			next(ctx)
			return
		}
		for _, scope := range credential.Scopes {
			if slices.Contains(anyOfNeededScopes, scope) {
				next(ctx)
				return
			}
		}
		huma.WriteErr(api, ctx, http.StatusForbidden, "Bearer token is not authorized for any of the scopes: "+strings.Join(anyOfNeededScopes, ", "))
	}
}

func authenticate(displayNameApp, email string) error {
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.IntN(10))
	}
	userName := email[:strings.Index(email, "@")]
	user, err := stores.NewUserStore(relationalDB).GetByEmail(email)
	if err != nil {
		return err
	}
	if user != nil {
		userName = user.Name
	}
	if err := mailer.SendCode(displayNameApp, userName, email, code); err != nil {
		return err
	}
	authenticationRequest := AuthenticationRequest{
		appName:  displayNameApp,
		userName: userName,
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
	if _, err := uuid.Parse(token); err != nil {
		return nil
	}
	credential, err := stores.NewCredentialStore(relationalDB).Get(token)
	if err != nil {
		log.Println("ERROR getting credential from store:", err)
	}
	if credential == nil {
		return nil
	}
	sessions.SetDefault(token, credential)
	return credential
}

func flushSession(userID string) {
	selectedKey := ""
	for key, item := range sessions.Items() {
		if item.Object.(*models.Credential).UserID == userID {
			selectedKey = key
			break
		}
	}
	if selectedKey != "" {
		sessions.Delete(selectedKey)
	}
}

func assertAdminUserExists(config *Configuration) error {
	roles, err := stores.NewRoleStore(relationalDB).GetAll()
	if err != nil {
		return err
	}
	var admin models.Role
	for _, r := range roles {
		if r.Name == "administrator" {
			admin = r
			break
		}
	}
	if admin.ID <= 0 {
		return errors.New("no administrator role was found in the relational database")
	}
	user, err := stores.NewUserStore(relationalDB).GetByEmail(config.AdminEmailAddress)
	if err != nil {
		return err
	}
	if user == nil {
		defaultAdminUser := models.UserCreatedByAdmin{}
		defaultAdminUser.Name = "_Administrator_"
		defaultAdminUser.Email = config.AdminEmailAddress
		user, err = stores.NewUserStore(relationalDB).Add(&defaultAdminUser)
		if err != nil {
			return err
		}
	}
	profile, err := stores.NewProfileStore(relationalDB).Get(user.ID)
	if err != nil {
		return err
	}
	for _, r := range profile.Roles {
		if r.ID == admin.ID {
			return nil
		}
	}
	return stores.NewRoleStore(relationalDB).AddRoleToUser(profile.ID, admin.ID)
}
