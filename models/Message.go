package models

type MessageRecord struct {
	ID           string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID for this message."`
	Name         string  `json:"name" doc:"The name of this message."`
	Severity     int     `json:"severity" minimum:"1" maximum:"5" doc:"The severity for this message, where 1 means 'info' and 5 means 'urgent'."`
	Text         string  `json:"text" doc:"The text of this message."`
	ExperimentID string  `json:"experimentID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the experiment this message belongs to."`
	SpeciesID    *string `json:"speciesID,omitempty" format:"uuid" writeOnly:"true" doc:"The ID of the species this message is associated with."`
	AnswerID     *string `json:"answerID,omitempty" format:"uuid" writeOnly:"true" doc:"The ID of the answer this message is associated with."`
}

type Message struct {
	MessageRecord
	Experiment Experiment `json:"experiment" doc:"The experiment that this questionnaire belongs to."`
	Species    *Species   `json:"species" doc:"The optional species this message is associated with."`
	Answer     *Answer    `json:"answer" doc:"The optional answer this message is associated with."`
}
