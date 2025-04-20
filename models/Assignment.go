package models

type Assignment struct {
	User          User          `json:"user" doc:"The user this assignment is for."`
	Questionnaire Questionnaire `json:"questionnaire" doc:"The questionnaire that was assigned."`
	Interaction   Interaction   `json:"interaction" doc:"The interaction that produced the assignment."`
}
