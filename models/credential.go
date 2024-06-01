package models

import "time"

type Credential struct {
	UserID    string    `json:"userID" format:"uuid" doc:"The ID of the user this credential belongs to."`
	Email     string    `json:"email" format:"email" doc:"The email address associated with this credential."`
	Token     string    `json:"token" format:"uuid" doc:"The bearer token associated with this credential."`
	Scopes    []string  `json:"scopes" doc:"The scopes that can be used with this credential."`
	LastLogon time.Time `json:"lastLogon" format:"date-time" doc:"The last time this credential was used to log on."`
}
