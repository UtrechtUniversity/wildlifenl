package models

type Belonging struct {
	ID       string `json:"ID" format:"uuid" doc:"The ID of this belonging."`
	Name     string `json:"name" minLength:"2" doc:"The name of this belonging."`
	Category string `json:"category" minLength:"2" doc:"The category of this belonging."`
}
