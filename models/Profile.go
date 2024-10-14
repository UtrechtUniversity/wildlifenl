package models

import "time"

type ProfileRecord struct {
	User
	Email string `json:"email" format:"email" readOnly:"true" doc:"The email address of this user."`
}

type Profile struct {
	ProfileRecord
	Roles             []Role     `json:"roles,omitempty" doc:"The additional roles this user has."`
	Location          *Point     `json:"location,omitempty" readOnly:"true" doc:"The location of this user."`
	LocationTimestamp *time.Time `json:"locationTimestamp,omitempty" format:"date-time" readOnly:"true" doc:"The moment that this user's location was last updated."`
}
