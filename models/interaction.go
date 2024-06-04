package models

type InteractionRecord struct {
	ID          string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this interaction."`
	Description string  `json:"description" doc:"The description of this interaction."`
	Latitude    float64 `json:"latitude" minimum:"-89.999999" maximum:"89.999999" doc:"The latitude of the location associated with this interaction."`
	Longitude   float64 `json:"longitude" minimum:"-89.999999" maximum:"89.999999" doc:"The longitude of the location associated with this interaction."`
	SpeciesID   string  `json:"speciesID" format:"uuid" writeOnly:"true"`
}

type Interaction struct {
	InteractionRecord
	Species Species `json:"species" doc:"The species of the animal that this interaction was with."`
	User    User    `json:"user" doc:"The User that reported this interaction."`
}
