package models

import "time"

type AlarmRecord struct {
	ID            string    `json:"ID" format:"uuid" doc:"The ID of this alarm."`
	Timestamp     time.Time `json:"timestamp" format:"date-time" doc:"The moment this alarm was created."`
	ZoneID        string    `json:"-"`
	InteractionID *string   `json:"-"`
	DetectionID   *string   `json:"-"`
	AnimalID      *string   `json:"-"`
}

type Alarm struct {
	AlarmRecord
	Zone        Zone         `json:"zone" doc:"The zone for which this alarm is."`
	Interaction *Interaction `json:"interaction,omitempty" doc:"The optional interaction that initiated the creation of this alarm."`
	Detection   *Detection   `json:"detection,omitempty" doc:"The optional detection that initiated the creation of this alarm."`
	Animal      *Animal      `json:"animal,omitempty" doc:"The optional animal that initiated the creation of this alarm."`
	Conveyances []Conveyance `json:"conveyances" doc:"The conveyances that were created for this alarm."`
}
