package models

type Animal struct {
	ID      string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this animal."`
	Name    string  `json:"name" doc:"The name of this animal." example:"Flupke"`
	Species Species `json:"species" doc:"The species of this animal"`
}
