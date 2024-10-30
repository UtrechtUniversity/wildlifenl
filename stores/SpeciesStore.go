package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type SpeciesStore Store

func NewSpeciesStore(db *sql.DB) *SpeciesStore {
	s := SpeciesStore{
		relationalDB: db,
		query: `
		SELECT s."ID", s."name", s."commonName"
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
		var s models.Species
		if err := rows.Scan(&s.ID, &s.Name, &s.CommonName); err != nil {
			return nil, err
		}
		speciesX = append(speciesX, s)
	}
	return speciesX, nil
}

func (s *SpeciesStore) Get(speciesID string) (*models.Species, error) {
	query := s.query + `
		WHERE s."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, speciesID)
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
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *SpeciesStore) Add(species *models.Species) (*models.Species, error) {
	query := `
		INSERT INTO "species"("name", "commonName") VALUES($1, $2)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, species.Name, species.CommonName)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
