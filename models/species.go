package models

type Species struct {
	ID         string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this species."`
	Name       string  `json:"name" doc:"The Latin binomen of this species."`
	CommonName string  `json:"commonName" doc:"The common name of this species."`
	Category   *string `json:"category,omitempty" doc:"The animal category for this species."`
	Advice     *string `json:"advice,omitempty" doc:"The advice on how to interact with this species."`
	DidYouKnow *string `json:"didYouKnow,omitempty" doc:"A nice info text on this species."`
}
