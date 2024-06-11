package wildlifenl

type Configuration struct {
	Host                        string
	Port                        int
	DbHost                      string
	DbName                      string
	DbUser                      string
	DbPass                      string
	DbSSLmode                   string
	CacheSessionDurationMin     int
	CacheAuthRequestDurationMin int
}
