package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AnimalStore Store

func NewAnimals(relationalDB *sql.DB, timeseriesDB *Timeseries) *AnimalStore {
	s := AnimalStore{
		relationalDB: relationalDB,
		timeseriesDB: timeseriesDB,
		query: `
		SELECT a."id", a."name", a."location", s."id", s."name", s."commonNameNL", s."commonNameEN"
		FROM animal a
		LEFT JOIN species s ON s.id = a."speciesID"
		`,
	}
	return &s
}

func (s *AnimalStore) Get(id string) (*models.Animal, error) {
	query := s.query + `
		WHERE a."id" = $1
	`
	var animal models.Animal
	var species models.Species
	row := s.relationalDB.QueryRow(query, id)
	if err := row.Scan(&animal.ID, &animal.Name, &animal.Location, &species.ID, &species.Name, &species.CommonNameNL, &species.CommonNameEN); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	animal.Species = species
	return &animal, nil
}
