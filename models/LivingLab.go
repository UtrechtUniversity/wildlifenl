package models

type LivingLab struct {
	ID         string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this living lab."`
	Name       string  `json:"name" doc:"The name of this living lab."`
	Definition Polygon `json:"definition" doc:"The polygon describing the area of this living lab."`
}
