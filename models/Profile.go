package models

import "time"

type ProfileRecord struct {
	User
	Email       string  `json:"email" format:"email" readOnly:"true" doc:"The email address of this user."`
	DateOfBirth *string `json:"dateOfBirth,omitempty" format:"date" doc:"The date of birth of this user."`
	Gender      *string `json:"gender,omitempty" enum:"female,male,other" doc:"The gender of this user."`
	Postcode    *string `json:"postcode,omitempty" doc:"The postcode of this user."`
	Description *string `json:"description,omitempty" doc:"The description of this user."`
}

type Profile struct {
	ProfileRecord
	Roles             []Role     `json:"roles,omitempty" doc:"The additional roles this user has."`
	Location          *Point     `json:"location,omitempty" readOnly:"true" doc:"The location of this user."`
	LocationTimestamp *time.Time `json:"locationTimestamp,omitempty" format:"date-time" readOnly:"true" doc:"The moment that this user's location was last updated."`
}
