package stores

import (
	"database/sql"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type TrackingEventStore Store

func NewTrackingEventStore(db *sql.DB) *TrackingEventStore {
	s := TrackingEventStore{
		db: db,
		query: `
		SELECT e."createdAt", e."latitude", e."longitude", u."id", u."name"
		FROM "trackingEvent" e
		INNER JOIN "user" u ON u."id" = e."userID"
		`,
	}
	return &s
}

func (s *TrackingEventStore) process(rows *sql.Rows, err error) ([]models.TrackingEvent, error) {
	if err != nil {
		return nil, err
	}
	trackingEvents := make([]models.TrackingEvent, 0)
	for rows.Next() {
		var trackingEvent models.TrackingEvent
		var user models.User
		if err := rows.Scan(&trackingEvent.CreatedAt, &trackingEvent.Latitude, &trackingEvent.Longitude, &user.ID, &user.Name); err != nil {
			return nil, err
		}
		trackingEvent.User = user
		trackingEvents = append(trackingEvents, trackingEvent)
	}
	return trackingEvents, nil
}

func (s *TrackingEventStore) Get(createdAt time.Time, userID string) (*models.TrackingEvent, error) {
	query := s.query + `
		WHERE e."createdAt" = $1
		AND u."id" = $2
		`
	rows, err := s.db.Query(query, createdAt, userID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *TrackingEventStore) GetByUser(userID string) ([]models.TrackingEvent, error) {
	query := s.query + `
		WHERE u."id" = $1
		ORDER BY e."createdAt" DESC
		`
	rows, err := s.db.Query(query, userID)
	return s.process(rows, err)
}

func (s *TrackingEventStore) Add(userID string, trackingEvent *models.TrackingEventRecord) (*models.TrackingEvent, error) {
	query := `
		INSERT INTO "trackingEvent"("latitude", "longitude", "userID") VALUES($1, $2, $3)
		RETURNING "createdAt"
	`
	var createdAt time.Time
	row := s.db.QueryRow(query, trackingEvent.Latitude, trackingEvent.Longitude, userID)
	if err := row.Scan(&createdAt); err != nil {
		return nil, err
	}
	return s.Get(createdAt, userID)
}
