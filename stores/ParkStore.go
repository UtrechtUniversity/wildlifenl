package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ParkStore Store

func NewParkStore(db *sql.DB) *ParkStore {
	s := ParkStore{
		relationalDB: db,
		query: `
		SELECT p."id", p."name", p."definition"
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
		if err := rows.Scan(&park.ID, &park.Name, &park.Definition); err != nil {
			return nil, err
		}
		parks = append(parks, park)
	}
	return parks, nil
}

func (s *ParkStore) Get(parkID string) (*models.Park, error) {
	query := s.query + `
		WHERE p."id" = $1
		`
	rows, err := s.relationalDB.Query(query, parkID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ParkStore) GetAll() ([]models.Park, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *ParkStore) Add(park *models.Park) (*models.Park, error) {
	query := `
		INSERT INTO park("name", "definition") VALUES($1, $2)
		RETURNING "id"
	`
	var id string
	row := s.relationalDB.QueryRow(query, park.Name, park.Definition)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
