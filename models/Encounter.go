package models

import "time"

type Encounter struct {
	ID             string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this encounter."`
	Timestamp      time.Time `json:"timestamp" readOnly:"true" doc:"The moment this encounter happened."`
	UserLocation   Point     `json:"userLocation" doc:"The location of the user involved in this encounter at the moment the encounter happened."`
	AnimalLocation Point     `json:"animalLocation" doc:"The location of the animal involved in this encounter at the moment the encouter happened."`
	User           User      `json:"user" doc:"The user involved in this encounter."`
	Animal         Animal    `json:"animal" doc:"The animal involved in this encounter."`
}
