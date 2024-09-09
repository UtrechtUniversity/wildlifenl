package main

import (
	"log"
	"os"
	"strconv"

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
			log.Fatal("environment variable CACHE_SESSION_MINUTES must be empty for the default value of 30, or an unsigned 16-bit integer: 10 <= x <= 65535")
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
	m, err := migrate.New("file://../database", connStr)
	if err != nil {
		return err
	}
	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}
