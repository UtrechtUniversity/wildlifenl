package models

type User struct {
	ID                 string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this user."`
	Name               string `json:"name" minLength:"2" doc:"The display name of this user."`
	RecreationAppTandC bool   `json:"recreationAppTerms" doc:"Reports whether this user accepted the terms and conditions for the use of the recreation app."`
	ReportAppTandC     bool   `json:"reportAppTerms" doc:"Reports whether this user accepted the terms and conditions for the use of the report app."`
}

type UserCreatedByAdmin struct {
	User
	Email string `json:"email" format:"email" doc:"The email address of this user."`
}
