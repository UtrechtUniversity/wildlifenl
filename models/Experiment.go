package models

import "time"

type ExperimentRecord struct {
	ID          string     `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this experiment."`
	Name        string     `json:"name" doc:"The name of this experiment."`
	Description string     `json:"description" doc:"The description of this experiment."`
	Start       time.Time  `json:"start" format:"date-time" doc:"The moment this experiment started."`
	End         *time.Time `json:"end,omitempty" format:"date-time" doc:"The moment this experiment ended."`
	LivingLabID *string    `json:"livingLabID,omitempty" format:"uuid" writeOnly:"true"`
}

type Experiment struct {
	ExperimentRecord
	User      User       `json:"user" doc:"The User that created this experiment."`
	LivingLab *LivingLab `json:"livingLab,omitempty" doc:"The livingLab this experiment is bound to."`
}
