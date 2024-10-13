package stores

import (
	"database/sql"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/timeseries"
)

type AnimalStore Store

func NewAnimalStore(relationalDB *sql.DB, timeseriesDB *timeseries.Timeseries) *AnimalStore {
	s := AnimalStore{
		relationalDB: relationalDB,
		timeseriesDB: timeseriesDB,
		query: `
		SELECT a."ID", a."name", a."location", a."locationTimestamp", s."ID", s."name", s."commonNameNL", s."commonNameEN"
		FROM "animal" a
		LEFT JOIN "species" s ON s."ID" = a."speciesID"
		`,
	}
	return &s
}

func (s *AnimalStore) process(rows *sql.Rows, err error) ([]models.Animal, error) {
	if err != nil {
		return nil, err
	}
	animals := make([]models.Animal, 0)
	for rows.Next() {
		var a models.Animal
		var s models.Species
		if err := rows.Scan(&a.ID, &a.Name, &a.Location, &a.LocationTimestamp, &s.ID, &s.Name, &s.CommonNameNL, &s.CommonNameEN); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		a.Species = s
		animals = append(animals, a)
	}
	return animals, nil
}

func (s *AnimalStore) Get(animalID string) (*models.Animal, error) {
	query := s.query + `
		WHERE a."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, animalID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *AnimalStore) GetAll() ([]models.Animal, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *AnimalStore) Add(animal *models.AnimalRecord) (*models.Animal, error) {
	query := `
		INSERT INTO "animal"("name", "speciesID") VALUES($1, $2)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, animal.Name, animal.SpeciesID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *AnimalStore) UpdateLocation(sensorID string, location models.Point, timestamp time.Time) (*models.Animal, error) {
	query := `
		UPDATE animal
		SET "location" = $1, "locationTimestamp" = $2
		WHERE "ID" = (
			SELECT a."ID"
			FROM "borneSensorDeployment" d
			INNER JOIN "animal" a ON a."ID" = d."animalID" 
			WHERE d."sensorID" = $3
			AND (d."end" IS NULL OR d."end" > $2)
			AND d."start" < $2
		)
		RETURNING "ID"
	`
	var id *string
	row := s.relationalDB.QueryRow(query, location, timestamp, sensorID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(*id)
}
