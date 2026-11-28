package models

type QuestionnaireRecord struct {
	ID                string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this questionnaire."`
	Name              string  `json:"name" minLength:"5" doc:"The name of this questionnaire."`
	Identifier        *string `json:"identifier,omitempty" doc:"An optional questionnaire identifier for internal use."`
	NotifyAfter       *string `json:"notifyAfter,omitempty" format:"time" doc:"An optional timestamp after which the user should be notified of having received this questionnaire if not immediately."`
	ExperimentID      string  `json:"experimentID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the experiment to associate this questionnaire with."`
	InteractionTypeID int     `json:"interactionTypeID,omitempty" writeOnly:"true" minimum:"1" required:"true" doc:"The ID of the interaction type to associate this questionnaire with."`
	SpeciesID         *string `json:"speciesID,omitempty" format:"uuid" writeOnly:"true" doc:"The optional ID of the species this questionnaire is bound to. If set, the questionnaire is only assigned to interactions with this species."`
}

type Questionnaire struct {
	QuestionnaireRecord
	Experiment      Experiment      `json:"experiment" doc:"The experiment that this questionnaire belongs to."`
	InteractionType InteractionType `json:"interactionType" doc:"The type of interactions that this questionnaire is created for."`
	Species         *Species        `json:"species,omitempty" doc:"The species this questionnaire is bound to."`
	Questions       []Question      `json:"questions,omitempty" doc:"The questions of this questionnaire"`
}
