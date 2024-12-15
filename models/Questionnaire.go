package models

type QuestionnaireRecord struct {
	ID                string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this questionnaire."`
	Name              string  `json:"name" doc:"The name of this questionnaire."`
	Identifier        *string `json:"identifier,omitempty" doc:"An optional questionnaire identifier for internal use."`
	ExperimentID      string  `json:"experimentID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the experiment to associate this questionnaire with."`
	InteractionTypeID int     `json:"interactionTypeID,omitempty" writeOnly:"true" minimum:"1" required:"true" doc:"The ID of the interaction type to associate this questionnaire with."`
}

type Questionnaire struct {
	QuestionnaireRecord
	Experiment      Experiment      `json:"experiment" doc:"The experiment that this questionnaire belongs to."`
	InteractionType InteractionType `json:"interactionType" doc:"The type of interactions that this questionnaire is created for."`
	Questions       []Question      `json:"questions,omitempty" doc:"The questions of this questionnaire"`
}
