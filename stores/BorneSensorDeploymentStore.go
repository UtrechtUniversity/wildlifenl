package stores

import (
	"database/sql"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/timeseries"
)

type BorneSensorDeploymentStore Store

func NewBorneSensorDeploymentStore(relationalDB *sql.DB, timeseriesDB *timeseries.Timeseries) *BorneSensorDeploymentStore {
	s := BorneSensorDeploymentStore{
		relationalDB: relationalDB,
		timeseriesDB: timeseriesDB,
		query: `
			SELECT d."ID", d."sensorID", d."contactHardwareAddress", d."start", d."end", a."ID", a."name", a."location", s."ID", s."name", s."commonName", COALESCE(u."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(u."name", '')
			FROM "borneSensorDeployment" d
			INNER JOIN "animal" a ON a."ID" = d."animalID"
			INNER JOIN "species" s ON s."ID" = a."speciesID"
			LEFT JOIN "user" u ON u."ID" = d."userID"
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
		var d models.BorneSensorDeployment
		var u models.User
		if err := rows.Scan(&d.ID, &d.SensorID, &d.ContactHardwareAddress, &d.Start, &d.End, &d.Animal.ID, &d.Animal.Identifier, &d.Animal.Location, &d.Animal.Species.ID, &d.Animal.Species.Name, &d.Animal.Species.CommonName, &u.ID, &u.Name); err != nil {
			return nil, err
		}
		if u.ID != "00000000-0000-0000-0000-000000000000" {
			d.User = u
		}
		borneSensorDeployments = append(borneSensorDeployments, d)
	}
	return borneSensorDeployments, nil
}

func (s *BorneSensorDeploymentStore) Get(ID string) (*models.BorneSensorDeployment, error) {
	query := s.query + `
		WHERE d."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, ID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *BorneSensorDeploymentStore) GetMine(userID string) ([]models.BorneSensorDeployment, error) {
	query := s.query + `
		WHERE d."userID" = $1
	`
	rows, err := s.relationalDB.Query(query, userID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *BorneSensorDeploymentStore) GetMyActive(userID string) ([]models.BorneSensorDeployment, error) {
	query := s.query + `
		WHERE d."end" IS NULL
		AND d."userID" = $1
	`
	rows, err := s.relationalDB.Query(query, userID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *BorneSensorDeploymentStore) GetByAnimal(animalID string, start time.Time, end time.Time) ([]models.BorneSensorDeployment, error) {
	query := s.query + `
		WHERE ("end" IS NULL OR "end" >= $2)
		AND "start" <= $3
		AND d."animalID" = $1
	`
	rows, err := s.relationalDB.Query(query, animalID, start, end)
	results, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	borneSensorDeployments := make([]models.BorneSensorDeployment, 0)
	for _, r := range results {
		readings, err := NewBorneSensorReadingStore(s.relationalDB, s.timeseriesDB).GetAllBySensorID(r.SensorID, start, end)
		if err != nil {
			return nil, err
		}
		r.BorneSensorReadings = readings
		borneSensorDeployments = append(borneSensorDeployments, r)
	}
	return borneSensorDeployments, nil
}

func (s *BorneSensorDeploymentStore) End(userID string, deploymentID string, end time.Time) (*models.BorneSensorDeployment, error) {
	query := `
		UPDATE "borneSensorDeployment"
		SET "end" = $1
		WHERE "ID" = $2
		AND "userID" = $3
		AND "end" IS NULL
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, end, deploymentID, userID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}

func (s *BorneSensorDeploymentStore) Add(userID string, borneSensorDeployment *models.BorneSensorDeploymentRecord) (*models.BorneSensorDeployment, error) {
	_, err := s.End(userID, borneSensorDeployment.SensorID, borneSensorDeployment.Start)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO "borneSensorDeployment"("animalID", "sensorID", "userID", "contactHardwareAddress", "start") VALUES($1, $2, $3, $4, $5)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, borneSensorDeployment.AnimalID, userID, borneSensorDeployment.SensorID, borneSensorDeployment.ContactHardwareAddress, borneSensorDeployment.Start)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
