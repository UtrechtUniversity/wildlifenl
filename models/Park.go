package models

type Park struct {
	ID         string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this park."`
	Name       string  `json:"name" doc:"The name of this park."`
	Definition Polygon `json:"definition" doc:"The polygon describing the area of this park."`
}
