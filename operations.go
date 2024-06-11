package wildlifenl

import "database/sql"

type Operations struct {
	Database *sql.DB
	Endpoint string
}
