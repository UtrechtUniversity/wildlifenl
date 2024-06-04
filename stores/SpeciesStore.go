package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type SpeciesStore Store

func NewSpeciesStore(db *sql.DB) *SpeciesStore {
	s := SpeciesStore{
		db: db,
		query: `
		SELECT s."id", s."name", s."commonNameNL", s."commonNameEN"
		FROM "species" s
		`,
	}
	return &s
}

func (s *SpeciesStore) process(rows *sql.Rows, err error) ([]models.Species, error) {
	if err != nil {
		return nil, err
	}
	speciesX := make([]models.Species, 0)
	for rows.Next() {
		var species models.Species
		rows.Scan(&species.ID, &species.Name, &species.CommonNameNL, &species.CommonNameEN)
		speciesX = append(speciesX, species)
	}
	return speciesX, nil
}

func (s *SpeciesStore) Get(speciesID string) (*models.Species, error) {
	query := s.query + `
		WHERE s."id" = $1
		`
	rows, err := s.db.Query(query, speciesID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *SpeciesStore) GetAll() ([]models.Species, error) {
	query := s.query + `
		ORDER BY s."name"
	`
	rows, err := s.db.Query(query)
	return s.process(rows, err)
}
