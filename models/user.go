package models

type User struct {
	ID   string `json:"id" format:"uuid" readOnly:"true" doc:"The ID of this user."`
	Name string `json:"name" doc:"The display name of this user."`
}
