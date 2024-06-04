package models

type Notice struct {
	ID          string     `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this notice."`
	Type        NoticeType `json:"type" doc:"The NoticeType of this notice."`
	Description string     `json:"description" doc:"The description of this notice."`
	Latitude    float64    `json:"latitude" minimum:"-89.999999" maximum:"89.999999" doc:"The latitude of the location associated with this notice."`
	Longitude   float64    `json:"longitude" minimum:"-89.999999" maximum:"89.999999" doc:"The longitude of the location associated with this notice."`
	Reporter    User       `json:"reporter,omitempty" readOnly:"true" doc:"The User that reported this notice."` // BUG: readOnly does not work?
}
