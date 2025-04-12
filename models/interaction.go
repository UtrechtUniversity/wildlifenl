package models

import "time"

type InteractionRecord struct {
	ID          string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this interaction."`
	Description string    `json:"description" doc:"The description of this interaction."`
	SpeciesID   string    `json:"speciesID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the species involved in this interaction."`
	Location    Point     `json:"location" doc:"The place where the interaction happened."`
	Moment      time.Time `json:"moment" doc:"The moment the interaction happened."`
	TypeID      int       `json:"typeID,omitempty" minimum:"1" writeOnly:"true" required:"true" doc:"The ID of the interaction type for this interaction."`
}

type Interaction struct {
	InteractionRecord
	Timestamp     time.Time       `json:"timestamp" format:"date-time" doc:"The moment this interaction was reported."`
	Species       Species         `json:"species" doc:"The species of the animal that this interaction was with."`
	User          User            `json:"user" doc:"The User that reported this interaction."`
	Type          InteractionType `json:"type" doc:"The type of the interaction."`
	Questionnaire *Questionnaire  `json:"questionnaire,omitempty" doc:"The questionnaire the user should fill-out after having added the interaction."`
}
