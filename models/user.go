package models

type User struct {
	ID   string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this user."`
	Name string `json:"name" minLength:"2" doc:"The display name of this user."`
}

type UserCreatedByAdmin struct {
	User
	Email string `json:"email" format:"email" doc:"The email address of this user."`
}
