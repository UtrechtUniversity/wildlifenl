package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ParkStore Store

func NewParkStore(db *sql.DB) *ParkStore {
	s := ParkStore{
		db: db,
		query: `
		SELECT p."id", u."name"
		FROM "park" p
		`,
	}
	return &s
}

func (s *ParkStore) process(rows *sql.Rows, err error) ([]models.Park, error) {
	if err != nil {
		return nil, err
	}
	parks := make([]models.Park, 0)
	for rows.Next() {
		var park models.Park
		rows.Scan(&park.ID, &park.Name, &park.Definition)
		parks = append(parks, park)
	}
	return parks, nil
}

func (s *ParkStore) GetAll() ([]models.Park, error) {
	rows, err := s.db.Query(s.query)
	return s.process(rows, err)
}
