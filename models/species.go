package models

type Species struct {
	ID           string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this species."`
	Name         string `json:"name" doc:"The Latin binomen"`
	CommonNameNL string `json:"commonNameNL" doc:"The Dutch common name"`
	CommonNameEN string `json:"commonNameEN" doc:"The English common name"`
}
