package models

import "time"

type TrackingReadingRecord struct {
	UserID    string    `json:"userID" format:"uuid" readOnly:"true" doc:"The ID of the user this reading belongs to."`
	Timestamp time.Time `json:"timestamp" format:"date-time" doc:"The moment that this reading was done."`
	Location  Point     `json:"location" doc:"The value read by the location sensor."`
}

type TrackingReading struct {
	TrackingReadingRecord
	Conveyance *Conveyance `json:"conveyance,omitempty" doc:"The optional conveyance of a message that should be displayed to the user."`
}
