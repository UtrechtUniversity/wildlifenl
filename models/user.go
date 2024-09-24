package models

type User struct {
	ID        string     `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this user."`
	Name      string     `json:"name" doc:"The display name of this user."`
	LivingLab *LivingLab `json:"livingLab,omitempty" doc:"The livingLab this user works at."`
}
