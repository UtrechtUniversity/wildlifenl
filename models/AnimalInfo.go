package models

type AnimalInfo struct {
	Sex       string `json:"sex" enum:"female,male,other" doc:"The sex of the observed animal."`
	LifeStage string `json:"lifeStage" enum:"infant,adolescent,adult,unknown" doc:"The life stage of the observed animal."`
	Condition string `json:"condition" enum:"healthy,impaired,dead,other" doc:"The condition of the observed animal."`
}
