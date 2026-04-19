package models

type CollisionReport struct {
	InvolvedAnimals []AnimalInfo `json:"involvedAnimals" doc:"Information on the animals involved in this animal-vehicle-collision report."`
	EstimatedDamage int          `json:"estimatedDamage" minimum:"0" doc:"The estimated value in Euros (€) of the damage."`
	Severity        string       `json:"severity" enum:"unknown,low,medium,high" doc:"The severity of the animal-vehicle-collision."`
}
