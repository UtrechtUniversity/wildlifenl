package models

type Me struct {
	User
	Email string `json:"email" format:"email" doc:"The email address of this user."`
}
