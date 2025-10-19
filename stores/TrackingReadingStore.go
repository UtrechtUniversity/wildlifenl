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

type TrackingReadingStore Store

func NewTrackingReadingStore(relationalDB *sql.DB, timeseriesDB *timeseries.Timeseries) *TrackingReadingStore {
	s := TrackingReadingStore{
		relationalDB: relationalDB,
		timeseriesDB: timeseriesDB,
		query: `
			from(bucket: "humans")
			|> range(start: -5y)
			|> filter(fn: (r) => r._measurement == "human")
		`,
	}
	return &s
}

func (s *TrackingReadingStore) process(records *api.QueryTableResult) ([]models.TrackingReading, error) {
	readings := make(map[string]map[time.Time]*models.TrackingReading)
	for records.Next() {
		r := records.Record()
		userID, ok := r.Values()["userID"].(string)
		if !ok {
			continue
		}
		sensor, ok := readings[userID]
		if !ok {
			sensor = make(map[time.Time]*models.TrackingReading)
			readings[userID] = sensor
		}
		reading, ok := sensor[r.Time()]
		if !ok {
			reading = &models.TrackingReading{}
			sensor[r.Time()] = reading
		}
		switch r.Field() {
		case "latitude":
			reading.Location.Latitude = r.Value().(float64)
		case "longitude":
			reading.Location.Longitude = r.Value().(float64)
		}
	}
	if err := records.Err(); err != nil {
		return nil, err
	}
	results := make([]models.TrackingReading, 0)
	for userID, timedReading := range readings {
		for time, reading := range timedReading {
			reading.UserID = userID
			reading.Timestamp = time
			results = append(results, *reading)
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Timestamp.After(results[j].Timestamp) })
	return results, nil
}

func (s *TrackingReadingStore) Add(userID string, trackingReading *models.TrackingReadingRecord) (*models.TrackingReading, error) {
	fields := make(map[string]any)
	fields["latitude"] = trackingReading.Location.Latitude
	fields["longitude"] = trackingReading.Location.Longitude
	writer := s.timeseriesDB.Writer("humans")
	tags := map[string]string{
		"userID": userID,
	}
	point := write.NewPoint("human", tags, fields, trackingReading.Timestamp.Local())
	if err := writer.WritePoint(context.Background(), point); err != nil {
		return nil, err
	}
	_, err := NewUserStore(s.relationalDB).UpdateLocation(userID, trackingReading.Location, trackingReading.Timestamp)
	if err != nil {
		return nil, err
	}
	trackingReading.UserID = userID
	return &models.TrackingReading{TrackingReadingRecord: *trackingReading}, nil
}

func (s *TrackingReadingStore) GetForUser(userID string) ([]models.TrackingReading, error) {
	query := s.query + `
		|> filter(fn: (r) => r.userID == "` + userID + `")
	`
	reader := s.timeseriesDB.Reader()
	records, err := reader.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer records.Close()
	return s.process(records)
}
