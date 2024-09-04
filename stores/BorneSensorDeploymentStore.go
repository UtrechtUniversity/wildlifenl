package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type BorneSensorDeploymentStore Store

func NewBorneSensorDeploymentStore(relationalDB *sql.DB) *BorneSensorDeploymentStore {
	s := BorneSensorDeploymentStore{
		relationalDB: relationalDB,
		query: `
			SELECT d."sensorID", d."start", d."end", a."ID", a."name", a."location", s."ID", s."name", s."commonNameNL", s."commonNameEN"
			FROM "borneSensorDeployment" d
			INNER JOIN "animal" a ON a."ID" = d."animalID"
			INNER JOIN "species" s ON s."ID" = a."speciesID"
			WHERE d."end" IS NULL
		`,
	}
	return &s
}

func (s *BorneSensorDeploymentStore) process(rows *sql.Rows, err error) ([]models.BorneSensorDeployment, error) {
	if err != nil {
		return nil, err
	}
	borneSensorDeployments := make([]models.BorneSensorDeployment, 0)
	for rows.Next() {
		var borneSensorDeployment models.BorneSensorDeployment
		var animal models.Animal
		var species models.Species
		if err := rows.Scan(&borneSensorDeployment.SensorID, &borneSensorDeployment.Start, &borneSensorDeployment.End, &animal.ID, &animal.Name, &animal.Location, &species.ID, &species.Name, &species.CommonNameNL, &species.CommonNameEN); err != nil {
			return nil, err
		}
		animal.Species = species
		borneSensorDeployment.Animal = animal
		borneSensorDeployments = append(borneSensorDeployments, borneSensorDeployment)
	}
	return borneSensorDeployments, nil
}

func (s *BorneSensorDeploymentStore) Get(sensorID string, animalID string) (*models.BorneSensorDeployment, error) {
	query := s.query + `
		AND d."sensorID" = $1
		AND d."animalID" = $2
	`
	rows, err := s.relationalDB.Query(query, sensorID, animalID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *BorneSensorDeploymentStore) GetAll() ([]models.BorneSensorDeployment, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *BorneSensorDeploymentStore) Add(borneSensorDeployment *models.BorneSensorDeploymentRecord) (*models.BorneSensorDeployment, error) {
	query := `
		UPDATE "borneSensorDeployment"
		SET "end" = $1
		WHERE "sensorID" = $2
	`
	if _, err := s.relationalDB.Exec(query, borneSensorDeployment.Start, borneSensorDeployment.SensorID); err != nil {
		return nil, err
	}
	query = `
		INSERT INTO "borneSensorDeployment"("animalID", "sensorID", "start") VALUES($1, $2, $3)
		RETURNING "sensorID", "animalID"
	`
	var sensorID string
	var animalID string
	row := s.relationalDB.QueryRow(query, borneSensorDeployment.AnimalID, borneSensorDeployment.SensorID, borneSensorDeployment.Start)
	if err := row.Scan(&sensorID, &animalID); err != nil {
		return nil, err
	}
	return s.Get(sensorID, animalID)
}
