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
        	|> range(start: -360d)
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
	return &models.TrackingReading{TrackingReadingRecord: *trackingReading}, nil
}
