package models

type NoticeType struct {
	ID     int    `json:"id" minimum:"1" doc:"The ID of this noticeType."`
	NameNL string `json:"nameNL" readOnly:"true" doc:"The Dutch name of this noticeType."`
	NameEN string `json:"nameEN" readOnly:"true" doc:"The English name of this noticeType."`
}
