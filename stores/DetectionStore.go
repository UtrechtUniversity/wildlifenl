package stores

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type DetectionStore Store

func NewDetectionStore(db *sql.DB) *DetectionStore {
	s := DetectionStore{
		relationalDB: db,
		query: `
		SELECT d."ID", d."location", d."timestamp", d."sensorID", d."sensorType", d."uri", s."ID", s."name", s."commonName"
		FROM "detection" d
		INNER JOIN "species" s ON s."ID" = d."speciesID"
		`,
	}
	return &s
}

func (s *DetectionStore) process(rows *sql.Rows, err error) ([]models.Detection, error) {
	if err != nil {
		return nil, err
	}
	detections := make([]models.Detection, 0)
	for rows.Next() {
		var d models.Detection
		if err := rows.Scan(&d.ID, &d.Location, &d.Timestamp, &d.SensorID, &d.SensorType, &d.URI, &d.Species.ID, &d.Species.Name, &d.Species.CommonName); err != nil {
			return nil, err
		}
		a, err := NewDetectionAnimalStore(s.relationalDB).GetAllForDetection(d.ID) // Potential performance issue, as it is being called inside a loop.
		if err != nil {
			return nil, err
		}
		d.Animals = a
		detections = append(detections, d)
	}
	return detections, nil
}

func (s *DetectionStore) Get(detectionID string) (*models.Detection, error) {
	query := s.query + `
		WHERE d."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, detectionID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *DetectionStore) GetAll() ([]models.Detection, error) {
	query := s.query + `
		ORDER BY d."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *DetectionStore) Add(detection models.DetectionRecord) (*models.Detection, error) {
	query := `
	INSERT INTO "detection"("location", "timestamp", "sensorID", "sensorType", "uri", "speciesID") VALUES($1, $2, $3, $4, $5, $6)
	RETURNING "ID"
`
	var id string
	row := s.relationalDB.QueryRow(query, detection.Location, detection.Timestamp, detection.SensorID, detection.SensorType, detection.URI, detection.SpeciesID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	if err := NewDetectionAnimalStore(s.relationalDB).addMany(detection.Animals, id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *DetectionStore) GetFiltered(area *models.Circle, before *time.Time, after *time.Time) ([]models.Detection, error) {
	query := s.query
	args := make([]any, 0)
	whereDone := false
	if area != nil {
		and := " AND "
		if !whereDone {
			and = " WHERE "
			whereDone = true
		}
		query += and + `$` + strconv.Itoa(len(args)+1) + `::circle @> d."location"`
		args = append(args, area)
	}
	if before != nil {
		and := " AND "
		if !whereDone {
			and = " WHERE "
			whereDone = true
		}
		query += and + `d."timestamp" < $` + strconv.Itoa(len(args)+1)
		args = append(args, before)
	}
	if after != nil {
		and := " AND "
		if !whereDone {
			and = " WHERE "
			whereDone = true
		}
		query += and + `d."timestamp" > $` + strconv.Itoa(len(args)+1)
		args = append(args, after)
	}
	rows, err := s.relationalDB.Query(query, args...)
	return s.process(rows, err)
}
