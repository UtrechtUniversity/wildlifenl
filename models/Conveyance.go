package models

import "time"

type Conveyance struct {
	ID        string     `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this conveyance."`
	Timestamp time.Time  `json:"timestamp" doc:"The moment this conveyance was created."`
	Message   Message    `json:"message" doc:"The message associated with this conveyance."`
	Encounter *Encounter `json:"encounter" doc:"The optional encounter associated with this conveyance."`
	Response  *Response  `json:"response" doc:"The optional response associated with this conveyance."`
}
