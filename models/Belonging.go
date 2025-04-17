package models

type Belonging struct {
	ID       string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this belonging."`
	Name     string `json:"name" doc:"The name of this belonging."`
	Category string `json:"category" doc:"The category of this belonging."`
}
