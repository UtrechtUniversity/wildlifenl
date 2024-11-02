package models

type MessageRecord struct {
	ID               string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID for this message."`
	Name             string  `json:"name" doc:"The name of this message."`
	Severity         int     `json:"severity" minimum:"1" maximum:"5" doc:"The severity for this message, where 1:debug 2:info 3:warning 4:urgent 5:critical."`
	Text             string  `json:"text" doc:"The text of this message."`
	ExperimentID     string  `json:"experimentID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the experiment this message belongs to."`
	Trigger          string  `json:"trigger" enum:"encounter,answer,alarm" doc:"The trigger type for this message."`
	SpeciesID        *string `json:"speciesID,omitempty" format:"uuid" writeOnly:"true" doc:"The optional ID of the species this message is associated with."`
	AnswerID         *string `json:"answerID,omitempty" format:"uuid" writeOnly:"true" doc:"The optional ID of the answer this message is associated with."`
	EncounterMeters  *int    `json:"encounterMeters,omitempty" minimum:"1" doc:"The distance in meters between location measurements of humans and the specified species for which to send this message. Only used if trigger = encounter."`
	EncounterMinutes *int    `json:"encounterMinutes,omitempty" minimum:"1" doc:"The time difference in minutes between location measurements of humans and the specified species. Only used if trigger = encounter."`
}

type Message struct {
	MessageRecord
	Experiment Experiment `json:"experiment" doc:"The experiment that this questionnaire belongs to."`
	Species    *Species   `json:"species,omitempty" doc:"The optional species this message is associated with."`
	Answer     *Answer    `json:"answer,omitempty" doc:"The optional answer this message is associated with."`
}
