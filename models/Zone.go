package models

type ZoneRecord struct {
	ID          string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this Zone."`
	Name        string `json:"name" doc:"The name of this Zone."`
	Description string `json:"description" doc:"A description for this Zone."`
	Area        Circle `json:"area" doc:"The geographic circle that defines this Zone."`
}

type Zone struct {
	ZoneRecord
	User    User      `json:"user" doc:"The user that added this Zone."`
	Species []Species `json:"species,omitempty" doc:"The species for which this zone creates alarms."`
}
