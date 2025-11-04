package models

type Vicinity struct {
	Animals      []Animal      `json:"animals" doc:"The animals in the vicinity."`
	Detections   []Detection   `json:"detections" doc:"The detections in the vicinity."`
	Interactions []Interaction `json:"interactions" doc:"The interactions in the vicinity."`
}
