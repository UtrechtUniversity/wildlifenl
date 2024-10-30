package models

type Species struct {
	ID         string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this species."`
	Name       string `json:"name" doc:"The Latin binomen of this species."`
	CommonName string `json:"commonName" doc:"The common name of this species."`
}
