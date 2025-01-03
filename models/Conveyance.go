package models

import "time"

type Conveyance struct {
	ID        string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this conveyance."`
	Timestamp time.Time `json:"timestamp" format:"date-time" doc:"The moment this conveyance was created."`
	User      User      `json:"user" doc:"The user this conveyance is for."`
	Message   Message   `json:"message" doc:"The message associated with this conveyance."`
	Animal    *Animal   `json:"animal,omitempty" doc:"The encountered animal, in case this conveyance was created for an encounter."`
	Response  *Response `json:"response,omitempty" doc:"The given response, in case this conveyance was created for a response."`
	Alarm     *Alarm    `json:"alarm,omitempty" doc:"The alarm, in case this conveyance was created for an alarm."`
}
