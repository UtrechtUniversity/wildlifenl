package models

import "time"

type InteractionRecord struct {
	ID                string           `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this interaction."`
	Description       string           `json:"description" doc:"The description of this interaction."`
	SpeciesID         string           `json:"speciesID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the species involved in this interaction."`
	Location          Point            `json:"location" doc:"The location where this interaction was reported."`
	Moment            time.Time        `json:"moment" format:"date-time" doc:"The moment this interaction happened."`
	Place             Point            `json:"place" doc:"The place where this interaction happened."`
	TypeID            int              `json:"typeID,omitempty" minimum:"1" writeOnly:"true" required:"true" doc:"The ID of the interaction type for this interaction."`
	ReportOfSighting  *SightingReport  `json:"reportOfSighting,omitempty" doc:"Report of the animal sightings. Only used for interactions of TypeID 1"`
	ReportOfDamage    *DamageReport    `json:"reportOfDamage,omitempty" doc:"Report of the inflicted damage. Only used for interactions of TypeID 2"`
	ReportOfCollision *CollisionReport `json:"reportOfCollision,omitempty" doc:"Report of the animal-vehicle-collision. Only used for interactions of TypeID 3"`
}

type Interaction struct {
	InteractionRecord
	Timestamp     time.Time       `json:"timestamp" format:"date-time" doc:"The date+time this interaction was reported."`
	Species       Species         `json:"species" doc:"The species of the animal that this interaction was with."`
	User          User            `json:"user" doc:"The User that reported this interaction."`
	Type          InteractionType `json:"type" doc:"The type of this interaction."`
	Questionnaire *Questionnaire  `json:"questionnaire,omitempty" doc:"The questionnaire the user should fill-out after having added this interaction."`
}
