package models

type AnswerRecord struct {
	ID             string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this answer."`
	Text           string  `json:"text" minLength:"5" doc:"The text of this answer."`
	Index          int     `json:"index" minimum:"1" doc:"The index of this answer within the question. If multiple answers have the same index the client application should present them to the end user in a randomly shuffled order."`
	QuestionID     string  `json:"questionID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the question this answer belongs to."`
	NextQuestionID *string `json:"nextQuestionID,omitempty" format:"uuid" doc:"The optional ID of the question in the same questionnaire that should follow after the user selected this answer."`
}

type Answer struct {
	AnswerRecord
}
