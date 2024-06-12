package models

type AreaRecord struct {
	ID          string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this area."`
	Description string  `json:"description" doc:"The description of this area"`
	Definition  Polygon `json:"definition" doc:"The polygon describing this area"`
}

type Area struct {
	AreaRecord
	User User `json:"user"`
}
