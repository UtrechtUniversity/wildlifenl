package stores

import "database/sql"

type Store struct {
	relationalDB *sql.DB
	timeseriesDB *Timeseries
	query        string
}
