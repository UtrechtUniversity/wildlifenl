package models

type CollisionReport struct {
	InvolvedAnimals []AnimalInfo `json:"involvedAnimals" doc:"Information on the animals involved in this animal-vehicle-collision report."`
	EstimatedDamage int          `json:"estimatedDamage" minimum:"0" doc:"The estimated value in Euros (€) of the damage."`
	Intensity       string       `json:"intensity" enum:"high,medium,low" doc:"The intensity of the animal-vehicle-collision."`
	Urgency         string       `json:"urgency" enum:"high,medium,low" doc:"The urgency of the animal-vehicle-collision."`
}
