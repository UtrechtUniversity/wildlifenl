package models

import "time"

type TrackingEventRecord struct {
	CreatedAt time.Time `json:"createdAt" readOnly:"true" doc:"The timestamp of the moment this tracking event was created."`
	Latitude  float64   `json:"latitude" minimum:"-89.999999" maximum:"89.999999" doc:"The latitude of the location associated with this tracking event."`
	Longitude float64   `json:"longitude" minimum:"-89.999999" maximum:"89.999999" doc:"The longitude of the location associated with this tracking event."`
}

type TrackingEvent struct {
	TrackingEventRecord
	User User `json:"user" doc:"The User that this tracking event is for."`
}
