package models

type DamageReportRecord struct {
	InteractionID   string `json:"interactionID" format:"uuid" doc:"The ID of the interaction this report belongs to."`
	BelongingID     string `json:"belongingID" format:"uuid" writeOnly:"true" doc:"The ID of the belonging that was damaged."`
	ImpactType      string `json:"impactType" enum:"square-meters,units" doc:"The type of the impact of this damage report."`
	ImpactValue     int    `json:"impactValue" minimum:"1" doc:"The value of the impact of this damage report."`
	EstimatedDamage int    `json:"estimatedDamage" minimum:"0" doc:"The estimated value in Euros (€) of the damage."`
	EstimatedLoss   int    `json:"estimatedLoss" minimum:"0" doc:"The estimated economical loss in Euros (€) as incurred by the damage."`
}

type DamageReport struct {
	DamageReportRecord
	Belonging Belonging `json:"belonging" doc:"The belonging that was damaged."`
}
