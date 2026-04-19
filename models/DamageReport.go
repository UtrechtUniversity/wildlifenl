package models

type DamageReport struct {
	Belonging                     string `json:"belonging" doc:"The belonging that was damaged."`
	PerceivedLoss                 string `json:"estimatedLoss" enum:"unknown,0-250,250-500,500-1000,1000-2000,2000-5000,5000+" doc:"The perceived economical loss in Euros (€) as incurred by the damage."`
	PreventiveMeasures            bool   `json:"preventiveMeasures" doc:"Whether preventive measures were installed or not."`
	PreventiveMeasuresDescription string `json:"preventiveMeasuresDescription" doc:"The description of the installed preventive measures."`
}
