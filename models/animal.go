package models

import "time"

type AnimalRecord struct {
	ID        string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this animal."`
	Name      string `json:"name" minLength:"2" doc:"The name of this animal."`
	SpeciesID string `json:"speciesID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the species of this animal."`
}

type Animal struct {
	AnimalRecord
	Species           Species    `json:"species" doc:"The species of this animal"`
	Location          *Point     `json:"location,omitempty" readOnly:"true" doc:"The location of this animal"`
	LocationTimestamp *time.Time `json:"locationTimestamp,omitempty" format:"date-time" readOnly:"true" doc:"The moment that this animal's location was last updated."`
}
