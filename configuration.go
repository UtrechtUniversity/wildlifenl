package wildlifenl

import "time"

type Configuration struct {
	Version                         string
	Host                            string
	Port                            int
	RelationalDatabaseHost          string
	RelationalDatabasePort          int
	RelationalDatabaseName          string
	RelationalDatabaseUser          string
	RelationalDatabasePass          string
	RelationalDatabaseSSLmode       string
	TimeseriesDatabaseURL           string
	TimeseriesDatabaseOrganization  string
	TimeseriesDatabaseToken         string
	CacheSessionDurationMinutes     int
	CacheAuthRequestDurationMinutes int
	EmailFrom                       string
	EmailHost                       string
	EmailUser                       string
	EmailPassword                   string
	AdminEmailAddress               string
	RateLimitIPWindow               time.Duration
	RateLimitIPMaxRequests          int
	RateLimitWindow                 time.Duration
	RateLimitMaxRequests            int
	RateLimitBaseCooldown           time.Duration
	RateLimitMaxCooldown            time.Duration
	RateLimitPenaltyDecay           time.Duration
	FirebaseCredentials             []byte
}
