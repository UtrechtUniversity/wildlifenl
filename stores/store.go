package stores

import "database/sql"

type Store struct {
	db    *sql.DB
	query string
}
