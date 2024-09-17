package models

type AnswerRecord struct {
	ID         string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this answer."`
	Text       string `json:"text" doc:"The text of this answer."`
	Index      int    `json:"index" doc:"The index of this answer within the question."`
	QuestionID string `json:"questionID,omitempty" format:"uuid" writeOnly:"true"`
}

type Answer struct {
	AnswerRecord
}
