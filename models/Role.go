package models

type Role struct {
	ID   int    `json:"ID" minimum:"1" doc:"The ID of this role."`
	Name string `json:"name" minLength:"2" doc:"The name of this role"`
}
