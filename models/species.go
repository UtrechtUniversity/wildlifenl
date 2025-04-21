package models

type Species struct {
	ID           string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this species."`
	Name         string `json:"name" minLength:"2" doc:"The Latin binomen of this species."`
	CommonName   string `json:"commonName" minLength:"2" doc:"The common name of this species."`
	Category     string `json:"category" doc:"The animal category for this species."`
	Advice       string `json:"advice" doc:"The advice on how to interact with this species."`
	RoleInNature string `json:"roleInNature" doc:"The role this species has in nature."`
	Description  string `json:"description" doc:"Information about this species."`
	Behaviour    string `json:"behaviour" doc:"Information on the behaviour of this species."`
}
