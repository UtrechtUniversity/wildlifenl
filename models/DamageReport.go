package models

type DamageReport struct {
	Belonging       Belonging `json:"belonging" doc:"The belonging that was damaged."`
	ImpactType      string    `json:"impactType" enum:"square-meters,units" doc:"The type of the impact of this damage report."`
	ImpactValue     int       `json:"impactValue" minimum:"1" doc:"The value of the impact of this damage report."`
	EstimatedDamage int       `json:"estimatedDamage" minimum:"0" doc:"The estimated value in Euros (€) of the damage."`
	EstimatedLoss   int       `json:"estimatedLoss" minimum:"0" doc:"The estimated economical loss in Euros (€) as incurred by the damage."`
}
