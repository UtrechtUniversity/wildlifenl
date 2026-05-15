package models

type SightingReport struct {
	InvolvedAnimals              []AnimalInfo `json:"involvedAnimals" doc:"Information on the animals involved in this sighting report."`
	HumanActivity                string       `json:"humanActivity" enum:"unknown,walking,cycling,mountain biking,walking the dog,horse riding,photography,relaxing,other..." doc:"The activity of the human as reported by the user."`
	HumanActivityOther           *string      `json:"humanActivityOther,omitempty" doc:"Specification of other human activity."`
	PerceivedAnimalActivity      string       `json:"perceivedAnimalActivity" enum:"unknown,walking,eating,looking around,fleeing,resting,other..." doc:"The activity of the animal(s) as reported by the user."`
	PerceivedAnimalActivityOther *string      `json:"perceivedAnimalActivityOther,omitempty" doc:"Specification of other animal activity."`
}
