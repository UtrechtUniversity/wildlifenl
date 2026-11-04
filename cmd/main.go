package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/UtrechtUniversity/wildlifenl"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var version string

func main() {
	config := new(wildlifenl.Configuration)

	config.Version = version

	config.Host = os.Getenv("HOST")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("environment variable PORT cannot be empty")
	}
	httpPort, err := strconv.Atoi(port)
	if err != nil || httpPort <= 1024 || httpPort > 65535 {
		log.Fatal("environment variable PORT must be an unsigned 16-bit integer: 1024 < x <= 65535")
	}
	config.Port = httpPort

	config.RelationalDatabaseHost = os.Getenv("POSTGRESQL_HOST")
	if config.RelationalDatabaseHost == "" {
		log.Fatal("environment variable POSTGRESQL_HOST cannot be empty")
	}
	relationalDatabasePort := os.Getenv("POSTGRESQL_PORT")
	if relationalDatabasePort != "" {
		port, err := strconv.Atoi(relationalDatabasePort)
		if err != nil || port <= 1024 || port > 65535 {
			log.Fatal("environment variable POSTGRESQL_PORT must be empty or an unsigned 16-bit integer: 1024 < x <= 65535")
		}
		config.RelationalDatabasePort = port
	}
	config.RelationalDatabaseName = os.Getenv("POSTGRESQL_DBNAME")
	if config.RelationalDatabaseName == "" {
		log.Fatal("environment variable POSTGRESQL_DBNAME cannot be empty")
	}
	config.RelationalDatabaseUser = os.Getenv("POSTGRESQL_USER")
	if config.RelationalDatabaseUser == "" {
		log.Fatal("environment variable POSTGRESQL_USER cannot be empty")
	}
	config.RelationalDatabasePass = os.Getenv("POSTGRESQL_PASSWORD")
	if config.RelationalDatabasePass == "" {
		log.Fatal("environment variable POSTGRESQL_PASSWORD cannot be empty")
	}
	config.RelationalDatabaseSSLmode = os.Getenv("POSTGRESQL_SSLMODE")
	if config.RelationalDatabaseSSLmode == "" {
		log.Fatal("environment variable POSTGRESQL_SSLMODE cannot be empty, use 'disable' instead")
	}

	config.TimeseriesDatabaseURL = os.Getenv("INFLUXDB_URL")
	if config.TimeseriesDatabaseURL == "" {
		log.Fatal("environment variable INFLUXDB_URL cannot be empty")
	}
	config.TimeseriesDatabaseOrganization = os.Getenv("INFLUXDB_ORGANIZATION")
	if config.TimeseriesDatabaseOrganization == "" {
		log.Fatal("environment variable INFLUXDB_ORGANIZATION cannot be empty")
	}
	config.TimeseriesDatabaseToken = os.Getenv("INFLUXDB_TOKEN")
	if config.TimeseriesDatabaseToken == "" {
		log.Fatal("environment variable INFLUXDB_TOKEN cannot be empty")
	}

	config.CacheSessionDurationMinutes = 120
	cacheSession := os.Getenv("CACHE_SESSION_MINUTES")
	if cacheSession != "" {
		cacheSessionMins, err := strconv.Atoi(cacheSession)
		if err != nil || cacheSessionMins <= 10 || cacheSessionMins > 65535 {
			log.Fatal("environment variable CACHE_SESSION_MINUTES must be empty for the default value of 120, or an unsigned 16-bit integer: 10 <= x <= 65535")
		}
		config.CacheSessionDurationMinutes = cacheSessionMins
	}

	config.CacheAuthRequestDurationMinutes = 30
	cacheAuthrequest := os.Getenv("CACHE_AUTHREQUEST_MINUTES")
	if cacheAuthrequest != "" {
		cacheAuthrequestMins, err := strconv.Atoi(cacheAuthrequest)
		if err != nil || cacheAuthrequestMins <= 10 || cacheAuthrequestMins > 65535 {
			log.Fatal("environment variable CACHE_AUTHREQUEST_MINUTES must be empty for the default value of 30, or an unsigned 16-bit integer: 10 <= x <= 65535")
		}
		config.CacheAuthRequestDurationMinutes = cacheAuthrequestMins
	}

	config.EmailFrom = os.Getenv("EMAIL_FROM")
	if config.EmailFrom == "" {
		log.Fatal("environment variable EMAIL_FROM cannot be empty")
	}
	config.EmailHost = os.Getenv("EMAIL_HOST")
	if config.EmailHost == "" {
		log.Fatal("environment variable EMAIL_HOST cannot be empty")
	}
	config.EmailUser = os.Getenv("EMAIL_USER")
	if config.EmailUser == "" {
		log.Fatal("environment variable EMAIL_USER cannot be empty")
	}
	config.EmailPassword = os.Getenv("EMAIL_PASS")
	if config.EmailPassword == "" {
		log.Fatal("environment variable EMAIL_PASS cannot be empty")
	}

	config.AdminEmailAddress = os.Getenv("ADMIN_EMAIL_ADDRESS")
	if config.AdminEmailAddress == "" {
		log.Fatal("environment variable ADMIN_EMAIL_ADDRESS cannot be empty")
	}

	config.RateLimitIPWindow = 1 * time.Minute
	rateLimitIPWindow := os.Getenv("RATE_LIMIT_IP_WINDOW_MINUTES")
	if rateLimitIPWindow != "" {
		rateLimitIPWindowMins, err := strconv.Atoi(rateLimitIPWindow)
		if err != nil || rateLimitIPWindowMins < 1 || rateLimitIPWindowMins > 10 {
			log.Fatal("environment variable RATE_LIMIT_IP_WINDOW_MINUTES must be empty or an integer: 1 <= x <= 10")
		}
		config.RateLimitIPWindow = (time.Duration(rateLimitIPWindowMins) * time.Minute)
	}

	config.RateLimitIPMaxRequests = 10
	rateLimitIPMaxRequests := os.Getenv("RATE_LIMIT_IP_MAX_REQUESTS")
	if rateLimitIPMaxRequests != "" {
		rateLimitIPMaxRequestsNum, err := strconv.Atoi(rateLimitIPMaxRequests)
		if err != nil || rateLimitIPMaxRequestsNum < 1 || rateLimitIPMaxRequestsNum > 300 {
			log.Fatal("environment variable RATE_LIMIT_IP_MAX_REQUESTS must be empty or an integer: 1 <= x <= 300")
		}
		config.RateLimitIPMaxRequests = rateLimitIPMaxRequestsNum
	}

	config.RateLimitWindow = 1 * time.Minute
	rateLimitWindow := os.Getenv("RATE_LIMIT_WINDOW_MINUTES")
	if rateLimitWindow != "" {
		rateLimitWindowMins, err := strconv.Atoi(rateLimitWindow)
		if err != nil || rateLimitWindowMins < 1 || rateLimitWindowMins > 10 {
			log.Fatal("environment variable RATE_LIMIT_WINDOW_MINUTES must be empty or an unsigned 16-bit integer: 1 <= x <= 10")
		}
		config.RateLimitWindow = (time.Duration(rateLimitWindowMins) * time.Minute)
	}

	config.RateLimitMaxRequests = 120
	rateLimitMaxRequests := os.Getenv("RATE_LIMIT_MAX_REQUESTS")
	if rateLimitMaxRequests != "" {
		rateLimitMaxRequestsNum, err := strconv.Atoi(rateLimitMaxRequests)
		if err != nil || rateLimitMaxRequestsNum < 60 || rateLimitMaxRequestsNum > 3600 {
			log.Fatal("environment variable RATE_LIMIT_MAX_REQUESTS must be empty or an unsigned 16-bit integer: 60 <= x <= 3600")
		}
		config.RateLimitMaxRequests = rateLimitMaxRequestsNum
	}

	config.RateLimitBaseCooldown = 15 * time.Second
	rateLimitBaseCooldown := os.Getenv("RATE_LIMIT_BASE_COOLDOWN_SECONDS")
	if rateLimitBaseCooldown != "" {
		rateLimitBaseCooldownSecs, err := strconv.Atoi(rateLimitBaseCooldown)
		if err != nil || rateLimitBaseCooldownSecs < 1 || rateLimitBaseCooldownSecs > 30 {
			log.Fatal("environment variable RATE_LIMIT_BASE_COOLDOWN_SECONDS must be empty or an unsigned 16-bit integer: 1 <= x <= 30")
		}
		config.RateLimitBaseCooldown = (time.Duration(rateLimitBaseCooldownSecs) * time.Second)
	}

	config.RateLimitMaxCooldown = 10 * time.Minute
	rateLimitMaxCooldown := os.Getenv("RATE_LIMIT_MAX_COOLDOWN_MINUTES")
	if rateLimitMaxCooldown != "" {
		rateLimitMaxCooldownMins, err := strconv.Atoi(rateLimitMaxCooldown)
		if err != nil || rateLimitMaxCooldownMins < 5 || rateLimitMaxCooldownMins > 60 {
			log.Fatal("environment variable RATE_LIMIT_MAX_COOLDOWN_MINUTES must be empty or an unsigned 16-bit integer: 5 <= x <= 60")
		}
		config.RateLimitMaxCooldown = (time.Duration(rateLimitMaxCooldownMins) * time.Minute)
	}

	config.RateLimitPenaltyDecay = 2 * time.Minute
	rateLimitPenaltyDecay := os.Getenv("RATE_LIMIT_PENALTY_DECAY_MINUTES")
	if rateLimitPenaltyDecay != "" {
		rateLimitPenaltyDecayMins, err := strconv.Atoi(rateLimitPenaltyDecay)
		if err != nil || rateLimitPenaltyDecayMins < 1 || rateLimitPenaltyDecayMins > 30 {
			log.Fatal("environment variable RATE_LIMIT_PENALTY_DECAY_MINUTES must be empty or an unsigned 16-bit integer: 1 <= x <= 30")
		}
		config.RateLimitPenaltyDecay = (time.Duration(rateLimitPenaltyDecayMins) * time.Minute)
	}

	firebaseCredentials := os.Getenv("FIREBASE_CREDENTIALS_JSON")
	if firebaseCredentials == "" {
		log.Fatal("environment variable FIREBASE_CREDENTIALS_JSON cannot be empty")
	}
	config.FirebaseCredentials = []byte(firebaseCredentials)

	if err := migrateDatabase(config); err != nil {
		log.Fatal("Database migration error:", err)
	}

	log.Fatal(wildlifenl.Start(config))
}

func migrateDatabase(config *wildlifenl.Configuration) error {
	connStr := "postgres://" + config.RelationalDatabaseUser + ":" + config.RelationalDatabasePass + "@" + config.RelationalDatabaseHost
	if config.RelationalDatabasePort > 0 {
		connStr += ":" + strconv.Itoa(config.RelationalDatabasePort)
	}
	connStr += "/" + config.RelationalDatabaseName + "?sslmode=" + config.RelationalDatabaseSSLmode
	m, err := migrate.New("file://../database/", connStr)
	if err != nil {
		return err
	}
	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}
