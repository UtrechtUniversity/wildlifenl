package models

type InteractionType struct {
	ID            int    `json:"ID" readOnly:"true" minimum:"1" doc:"The ID of this interaction type."`
	NameNL        string `json:"nameNL" doc:"The Dutch name of this interaction type."`
	NameEN        string `json:"nameEN" doc:"The English name of this interaction type."`
	DescriptionNL string `json:"descriptionNL" doc:"The Dutch description of this interaction type."`
	DescriptionEN string `json:"descriptionEN" doc:"The English description of this interaction type."`
}
