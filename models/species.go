package models

type Species struct {
	ID               string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this species."`
	Name             string `json:"name" doc:"The Latin binomen"`
	CommonNameNL     string `json:"commonNameNL" doc:"The Dutch common name"`
	CommonNameEN     string `json:"commonNameEN" doc:"The English common name"`
	EncounterMeters  int    `json:"encounterMeters" minimum:"1" doc:"The distance in meters within which encounters are registered for humans and this species."`
	EncounterMinutes int    `json:"encounterMinutes" minimum:"1" doc:"The time difference in minutes within which encounters are registered for humans and this species."`
}
