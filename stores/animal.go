package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type Animals struct {
	db *sql.DB
}

func NewAnimals(db *sql.DB) *Animals {
	return &Animals{db: db}
}

func (s *Animals) Get(id string) (*models.Animal, error) {
	query := `
			SELECT a."id", a."name", s."id", s."name", s."commonNameNL", s."commonNameEN"
			FROM animal a
			LEFT JOIN species s ON s.id = a."speciesID"
			WHERE a."id" = $1
		`
	var animal models.Animal
	var species models.Species
	row := s.db.QueryRow(query, id)
	if err := row.Scan(&animal.ID, &animal.Name, &species.ID, &species.Name, &species.CommonNameNL, &species.CommonNameEN); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	animal.Species = species
	return &animal, nil
}
