package models

type QuestionnaireRecord struct {
	ID                string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this questionnaire."`
	Name              string `json:"name" doc:"The name of this questionnaire."`
	ExperimentID      string `json:"experimentID" format:"uuid" writeOnly:"true" required:"true"`
	InteractionTypeID int    `json:"interactionTypeID" writeOnly:"true" required:"true"`
}

type Questionnaire struct {
	QuestionnaireRecord
	Experiment      Experiment      `json:"experiment" doc:"The experiment that this questionnaire belongs to."`
	InteractionType InteractionType `json:"interactionType" doc:"The type of interactions that this questionnaire is created for."`
	Questions       []Question      `json:"questions,omitempty" doc:"The questions of this questionnaire"`
}
