package models

type InteractionType struct {
	ID          int    `json:"ID" readOnly:"true" minimum:"1" doc:"The ID of this interaction type."`
	Name        string `json:"name" doc:"The name of this interaction type."`
	Description string `json:"description" doc:"The description of this interaction type."`
}
