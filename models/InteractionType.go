package models

type InteractionType struct {
	ID          int    `json:"ID" readOnly:"true" minimum:"1" doc:"The ID of this interaction type."`
	Name        string `json:"name" minLength:"2" doc:"The name of this interaction type."`
	Description string `json:"description" minLength:"5" doc:"The description of this interaction type."`
}
