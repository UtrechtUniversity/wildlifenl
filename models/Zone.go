package models

import "time"

type ZoneRecord struct {
	ID          string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this Zone."`
	Created     time.Time `json:"created" readOnly:"true" doc:"The moment this zone was created."`
	Name        string    `json:"name" minLength:"2" doc:"The name of this Zone."`
	Description string    `json:"description" minLength:"5" doc:"The description for this Zone."`
	Area        Circle    `json:"area" doc:"The geographic circle that defines this Zone."`
}

type Zone struct {
	ZoneRecord
	Deactivated *time.Time `json:"deactivated,omitempty" doc:"The moment this zone record was deactivated."`
	User        User       `json:"user" doc:"The user that added this Zone."`
	Species     []Species  `json:"species,omitempty" doc:"The species for which this zone creates alarms."`
}
