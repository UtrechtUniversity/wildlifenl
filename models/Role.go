package models

type Role struct {
	ID   int    `json:"ID" doc:"The ID of this role."`
	Name string `json:"name" doc:"The name of this role"`
}
