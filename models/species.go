package models

type Species struct {
	ID           string `json:"id" format:"uuid" readOnly:"true" doc:"The ID of this species."`
	Name         string `json:"name" doc:"The latin binomen" example:"Canis familiaris"`
	CommonNameNL string `json:"commonNameNL" doc:"The Dutch common name" example:"Hond" `
	CommonNameEN string `json:"commonNameEN" doc:"The English common name" example:"Dog" `
}
