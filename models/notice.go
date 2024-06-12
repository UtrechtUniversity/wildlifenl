package models

import "time"

type NoticeRecord struct {
	ID          string    `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this notice."`
	Timestamp   time.Time `json:"timestamp" readOnly:"true" doc:"The moment that this notice was created"`
	TypeID      int       `json:"typeID,omitempty" writeOnly:"true"`
	Description string    `json:"description" doc:"The description of this notice."`
	Location    Point     `json:"location" doc:"The location that this notice was created at."`
}

type Notice struct {
	NoticeRecord
	User User       `json:"user" doc:"The User that reported this notice."`
	Type NoticeType `json:"type" doc:"The NoticeType of this notice."`
}
