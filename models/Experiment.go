package models

import "time"

type ExperimentRecord struct {
	ID          string     `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this experiment."`
	Name        string     `json:"name" minLength:"2" doc:"The name of this experiment."`
	Description string     `json:"description" minLength:"5" doc:"The description of this experiment."`
	Start       time.Time  `json:"start" format:"date-time" doc:"The moment this experiment started."`
	End         *time.Time `json:"end,omitempty" format:"date-time" doc:"The moment this experiment ended."`
	LivingLabID *string    `json:"livingLabID,omitempty" format:"uuid" writeOnly:"true" doc:"The optional ID of the living lab this experiment is bound to."`
}

type Experiment struct {
	ExperimentRecord
	User                   User       `json:"user" doc:"The User that created this experiment."`
	LivingLab              *LivingLab `json:"livingLab,omitempty" doc:"The livingLab this experiment is bound to."`
	NumberOfQuestionnaires *int       `json:"numberOfQuestionnaires,omitempty" minimum:"0" doc:"The number of questionnaires associated with this experiment."`
	NumberOfMessages       *int       `json:"numberOfMessages,omitempty" minimum:"0" doc:"The number of messages associated with this experiment."`
	QuestionnaireActivity  *int       `json:"questionnaireActivity,omitempty" minimum:"0" doc:"The number of questionnaires with at least one response."`
	MessageActivity        *int       `json:"messageActivity,omitempty" minimum:"0" doc:"The number of conveyances."`
}
