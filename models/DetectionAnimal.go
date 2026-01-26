package models

type DetectionAnimal struct {
	Confidence  int     `json:"confidence" minimum:"0" maximum:"100" doc:"The confidence value of the animal detection converted to a percentage."`
	Behaviour   *string `json:"behaviour,omitempty" doc:"A textual description of the behaviour of the animal."`
	Description *string `json:"description,omitempty" doc:"A textual description of the animal detection."`
	Sex         *string `json:"sex,omitempty" enum:"female,male" doc:"The sex of the detected animal."`
	LifeStage   *string `json:"lifeStage,omitempty" enum:"infant,adolescent,adult" doc:"The life stage of the detected animal."`
	Condition   *string `json:"condition,omitempty" enum:"healthy,impaired,dead" doc:"The condition of the detected animal."`
}
