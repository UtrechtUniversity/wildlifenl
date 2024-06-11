package models

import "time"

type InteractionRecord struct {
	ID          string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this interaction."`
	CreatedAt   time.Time `json:"createdAt" format:"date-time" readOnly:"true" doc:"The timestamp of the moment this tracking event was created."`
	Description string    `json:"description" doc:"The description of this interaction."`
	Latitude    float64   `json:"latitude" minimum:"-89.99999" maximum:"89.99999" doc:"The latitude of the location associated with this interaction."`
	Longitude   float64   `json:"longitude" minimum:"-179.99999" maximum:"179.999999" doc:"The longitude of the location associated with this interaction."`
	SpeciesID   string    `json:"speciesID,omitempty" format:"uuid" writeOnly:"true"`
}

type Interaction struct {
	InteractionRecord
	Species Species `json:"species" doc:"The species of the animal that this interaction was with."`
	User    User    `json:"user" doc:"The User that reported this interaction."`
}
