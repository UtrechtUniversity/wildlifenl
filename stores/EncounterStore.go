package stores

import (
	"database/sql"
	"strings"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type EncounterStore Store

func NewEncounterStore(db *sql.DB) *EncounterStore {
	s := EncounterStore{
		relationalDB: db,
		query: `
		SELECT e."ID", e."timestamp", e."userLocation", e."animalLocation", u."ID", u."name", a."ID", a."name", s."ID", s."name", s."commonNameNL", s."commonNameEN", s."encounterMeters", s."encounterMinutes"
		FROM "encounter" e
		INNER JOIN "user" u ON u."ID" = e."userID"
		INNER JOIN "animal" a ON a."ID" = e."animalID"
		INNER JOIN "species" s ON s."ID" = a."speciesID"
		`,
	}
	return &s
}

func (s *EncounterStore) process(rows *sql.Rows, err error) ([]models.Encounter, error) {
	if err != nil {
		return nil, err
	}
	encounters := make([]models.Encounter, 0)
	for rows.Next() {
		var e models.Encounter
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.UserLocation, &e.AnimalLocation, &e.User.ID, &e.User.Name, &e.Animal.ID, &e.Animal.Name, &e.Animal.Species.ID, &e.Animal.Species.Name, &e.Animal.Species.CommonNameNL, &e.Animal.Species.CommonNameEN, &e.Animal.Species.EncounterMeters, &e.Animal.Species.EncounterMinutes); err != nil {
			return nil, err
		}
		encounters = append(encounters, e)
	}
	return encounters, nil
}

func (s *EncounterStore) Get(encounterID int) (*models.Encounter, error) {
	query := s.query + `
		WHERE e."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, encounterID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *EncounterStore) GetAll() ([]models.Encounter, error) {
	query := s.query + `
		ORDER BY e."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *EncounterStore) AddAllForTrackingReading(trackingReading *models.TrackingReading) ([]models.Encounter, error) {
	query := `
		WITH inserted AS (
			INSERT INTO "encounter"("timestamp", "userLocation", "animalLocation", "userID", "animalID") 
			SELECT $1, $2, a."location", $3, a."ID"
			FROM "animal" a
			INNER JOIN "species" s ON s."ID" = a."speciesID"
			WHERE extract(epoch from $1 - a."locationTimestamp") / 60 > s."encounterMinutes"
			AND CIRCLE(a."location", CAST(s."encounterMeters" AS FLOAT) / 10000) @> $2
		    RETURNING "ID", "timestamp", "userLocation", "animalLocation", "userID", "animalID"
		)
	` + strings.Replace(s.query, "FROM \"encounter\"", "FROM inserted", 1)
	rows, err := s.relationalDB.Query(query, trackingReading.Timestamp, trackingReading.Location, trackingReading.UserID)
	if err != nil {
		return nil, err
	}
	return s.process(rows, err)
}
