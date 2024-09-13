package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/timeseries"
)

type Store struct {
	relationalDB *sql.DB
	timeseriesDB *timeseries.Timeseries
	query        string
}
