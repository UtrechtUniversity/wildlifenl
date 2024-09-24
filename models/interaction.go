package models

import "time"

type InteractionRecord struct {
	ID          string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this interaction."`
	Timestamp   time.Time `json:"timestamp" format:"date-time" readOnly:"true" doc:"The moment this interaction was reported."`
	Description string    `json:"description" doc:"The description of this interaction."`
	SpeciesID   string    `json:"speciesID,omitempty" format:"uuid" writeOnly:"true"`
	Location    Point     `json:"location" doc:"The place where the interaction was reported."`
	TypeID      int       `json:"typeID,omitempty" minimum:"1" writeOnly:"true"`
}

type Interaction struct {
	InteractionRecord
	Species       Species         `json:"species" doc:"The species of the animal that this interaction was with."`
	User          User            `json:"user" doc:"The User that reported this interaction."`
	Type          InteractionType `json:"type" doc:"The type of the interaction."`
	Questionnaire *Questionnaire  `json:"questionnaire,omitempty" doc:"The questionnaire the user should fill-out after having added the interaction."`
}
