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

type CannotUpdateError struct {
	message string
}

func (e *CannotUpdateError) Error() string {
	return e.message
}
