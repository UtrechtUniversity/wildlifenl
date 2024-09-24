package models

import "time"

type Profile struct {
	User
	Email             string     `json:"email" format:"email" doc:"The email address of this user."`
	Roles             []Role     `json:"roles,omitempty" doc:"The additional roles this user has."`
	Location          *Point     `json:"location,omitempty" readOnly:"true" doc:"The location of this user."`
	LocationTimestamp *time.Time `json:"locationTimestamp,omitempty" format:"datetime" readOnly:"true" doc:"The moment that this user's location was last updated."`
}
