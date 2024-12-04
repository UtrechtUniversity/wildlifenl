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

type ErrRecordInattainable struct {
	message string
}

func (e *ErrRecordInattainable) Error() string {
	return e.message
}

type ErrRecordImmutable struct {
	message string
}

func (e *ErrRecordImmutable) Error() string {
	return e.message
}
