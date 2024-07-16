package stores

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type BorneSensorReadingStore Store

func NewBorneSensorReadingStore(relationalDB *sql.DB, timeseriesDB *Timeseries) *BorneSensorReadingStore {
	s := BorneSensorReadingStore{
		relationalDB: relationalDB,
		timeseriesDB: timeseriesDB,
		query: `
			from(bucket: "animals")
        	|> range(start: -60m)
		`,
	}
	return &s
}

func (s *BorneSensorReadingStore) GetAll() ([]models.BorneSensorReading, error) {
	reader := s.timeseriesDB.Reader()
	records, err := reader.Query(context.Background(), s.query)
	if err != nil {
		return nil, err
	}
	readings := make(map[string]map[time.Time]*models.BorneSensorReading)
	for records.Next() {
		r := records.Record()
		sensorID, ok := r.Values()["sensorID"].(string)
		if !ok {
			continue
		}
		sensor, ok := readings[sensorID]
		if !ok {
			sensor = make(map[time.Time]*models.BorneSensorReading)
			readings[sensorID] = sensor
		}
		reading, ok := sensor[r.Time()]
		if !ok {
			reading = &models.BorneSensorReading{Location: models.Point{}}
			sensor[r.Time()] = reading
		}
		switch r.Field() {
		case "latitude":
			reading.Location.Latitude = r.Value().(float64)
		case "longitude":
			reading.Location.Longitude = r.Value().(float64)
		case "heartbeat":
			if value, ok := r.Value().(int64); ok {
				v := int(value)
				reading.Heartbeat = &v
			}
		case "temperature":
			if value, ok := r.Value().(int64); ok {
				v := int(value)
				reading.Temperature = &v
			}
		default:
			fmt.Println("unknown field:", r.Field())
		}
	}
	if err := records.Err(); err != nil {
		return nil, err
	}
	results := make([]models.BorneSensorReading, 0)
	for sensorID, timedReading := range readings {
		for time, reading := range timedReading {
			reading.SensorID = sensorID
			reading.Timestamp = time
			results = append(results, *reading)
		}
	}
	return results, nil
}

func (s *BorneSensorReadingStore) Add(borneSensorReading *models.BorneSensorReading) error {
	writer := s.timeseriesDB.Writer("animal")
	tags := map[string]string{
		"sensorID": borneSensorReading.SensorID,
	}
	fields := make(map[string]any)
	fields["latitude"] = borneSensorReading.Location.Latitude
	fields["longitude"] = borneSensorReading.Location.Longitude
	if borneSensorReading.Heartbeat != nil {
		fields["heartbeat"] = *borneSensorReading.Heartbeat
	}
	if borneSensorReading.Temperature != nil {
		fields["temperature"] = *borneSensorReading.Temperature
	}
	point := write.NewPoint("borne-sensor", tags, fields, time.Now())
	if err := writer.WritePoint(context.Background(), point); err != nil {
		return err
	}
	query := `
		UPDATE animal SET "location" = $1, "locationTimestamp" = now()
		FROM animal a
		JOIN "borneSensorDeployment" d on a."id" = d."animalID"
		WHERE d."sensorID" = $2
	`
	_, err := s.relationalDB.Exec(query, borneSensorReading.Location, borneSensorReading.SensorID)
	if err != nil {
		return err
	}
	return nil
}
