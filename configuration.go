package wildlifenl

type Configuration struct {
	Host                            string
	Port                            int
	RelationalDatabaseHost          string
	RelationalDatabaseName          string
	RelationalDatabaseUser          string
	RelationalDatabasePass          string
	RelationalDatabaseSSLmode       string
	TimeseriesDatabaseURL           string
	TimeseriesDatabaseOrganization  string
	TimeseriesDatabaseToken         string
	CacheSessionDurationMinutes     int
	CacheAuthRequestDurationMinutes int
}
