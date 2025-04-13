package models

type SightingReport struct {
	InvolvedAnimals []AnimalInfo `json:"involvedAnimals" doc:"Information on the animals involved in this sighting report."`
}
