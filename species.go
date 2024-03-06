package wildlifenl

type Species struct {
	Name       string `json:"name" example:"Canis familiaris" doc:"The latin binomen"`
	CommonName string `json:"commonName" example:"Dog" doc:"The common name"`
}
