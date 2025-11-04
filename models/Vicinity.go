package models

type Vicinity struct {
	Animals      []Animal      `json:"animals,omitempty" doc:"The animals in the vicinity."`
	Detections   []Detection   `json:"detections,omitempty" doc:"The detections in the vicinity."`
	Interactions []Interaction `json:"interactions,omitempty" doc:"The interactions in the vicinity."`
}
