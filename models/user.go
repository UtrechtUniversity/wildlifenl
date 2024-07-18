package models

type User struct {
	ID    string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this user."`
	Name  string `json:"name" doc:"The display name of this user."`
	Roles []Role `json:"roles,omitempty" doc:"The additional roles this user has."`
}
