package models

import "time"

type TrackingReading struct {
	UserID    string    `json:"userID" format:"uuid" readOnly:"true" doc:"The ID of the user this reading belongs to."`
	Timestamp time.Time `json:"timestamp" format:"date-time" doc:"The moment that this reading was done."`
	Location  Point     `json:"location" doc:"The value read by the location sensor."`
}
