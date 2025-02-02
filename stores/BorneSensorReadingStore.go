package stores

import (
	"context"
	"database/sql"
	"sort"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/timeseries"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type BorneSensorReadingStore Store

func NewBorneSensorReadingStore(relationalDB *sql.DB, timeseriesDB *timeseries.Timeseries) *BorneSensorReadingStore {
	s := BorneSensorReadingStore{
		relationalDB: relationalDB,
		timeseriesDB: timeseriesDB,
		query: `
			from(bucket: "animals")
        	|> range(start: -360d)
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
	return s.process(records)
}

func (s *BorneSensorReadingStore) GetAllBySensorID(sensorID string) ([]models.BorneSensorReading, error) {
	query := s.query + `
		|> filter(fn: (r) => r["sensorID"] == "` + sensorID + `")
	`
	reader := s.timeseriesDB.Reader()
	records, err := reader.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return s.process(records)
}

func (s *BorneSensorReadingStore) process(records *api.QueryTableResult) ([]models.BorneSensorReading, error) {
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
			reading = &models.BorneSensorReading{}
			sensor[r.Time()] = reading
		}
		switch r.Field() {
		case "userID":
			if value, ok := r.Value().(string); ok {
				reading.UserID = value
			}
		case "latitude":
			if reading.Location == nil {
				reading.Location = &models.Point{}
			}
			reading.Location.Latitude = r.Value().(float64)
		case "longitude":
			if reading.Location == nil {
				reading.Location = &models.Point{}
			}
			reading.Location.Longitude = r.Value().(float64)
		case "altitude":
			if value, ok := r.Value().(float64); ok {
				reading.Altitude = &value
			}
		case "temperature":
			if value, ok := r.Value().(float64); ok {
				reading.Temperature = &value
			}
		case "acceleroX":
			if reading.Accelero == nil {
				reading.Accelero = &models.Accelero{}
			}
			reading.Accelero.X = r.Value().(float64)
		case "acceleroY":
			if reading.Accelero == nil {
				reading.Accelero = &models.Accelero{}
			}
			reading.Accelero.Y = r.Value().(float64)
		case "acceleroZ":
			if reading.Accelero == nil {
				reading.Accelero = &models.Accelero{}
			}
			reading.Accelero.Z = r.Value().(float64)
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
	sort.Slice(results, func(i, j int) bool { return results[i].Timestamp.After(results[j].Timestamp) })
	return results, nil
}

func (s *BorneSensorReadingStore) Add(userID string, borneSensorReading *models.BorneSensorReadingRecord) (*models.Animal, error) {
	fields := make(map[string]any)
	fields["userID"] = userID
	if borneSensorReading.Location != nil {
		fields["latitude"] = borneSensorReading.Location.Latitude
		fields["longitude"] = borneSensorReading.Location.Longitude
	}
	if borneSensorReading.Altitude != nil {
		fields["altitude"] = *borneSensorReading.Altitude
	}
	if borneSensorReading.Temperature != nil {
		fields["temperature"] = *borneSensorReading.Temperature
	}
	if borneSensorReading.Accelero != nil {
		fields["acceleroX"] = borneSensorReading.Accelero.X
		fields["acceleroY"] = borneSensorReading.Accelero.Y
		fields["acceleroZ"] = borneSensorReading.Accelero.Z
	}
	if len(fields) > 0 {
		writer := s.timeseriesDB.Writer("animals")
		tags := map[string]string{
			"sensorID": borneSensorReading.SensorID,
		}
		point := write.NewPoint("borne-sensor", tags, fields, borneSensorReading.Timestamp.Local())
		if err := writer.WritePoint(context.Background(), point); err != nil {
			return nil, err
		}
	}
	if borneSensorReading.Location != nil {
		animal, err := NewAnimalStore(s.relationalDB, s.timeseriesDB).UpdateLocation(borneSensorReading.SensorID, *borneSensorReading.Location, borneSensorReading.Timestamp)
		if err != nil {
			return nil, err
		}
		if animal != nil {
			return animal, nil
		}
	}
	return nil, nil
}
