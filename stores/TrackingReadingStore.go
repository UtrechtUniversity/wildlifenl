package stores

import (
	"context"
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/timeseries"
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
	result, err := reader.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	trackingReadings := make([]models.TrackingReading, 0)
	for result.Next() {
		record := result.Record()
		var trackingReading models.TrackingReading
		trackingReading.Timestamp = record.Time()
		trackingReading.UserID = record.ValueByKey("userID").(string)
		trackingReading.Location = models.Point{
			Latitude:  record.ValueByKey("latitude").(float64),
			Longitude: record.ValueByKey("longitude").(float64),
		}
		trackingReadings = append(trackingReadings, trackingReading)
	}
	if result.Err() != nil {
		return nil, err
	}
	return trackingReadings, nil
}
