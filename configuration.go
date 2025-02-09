package wildlifenl

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
}
