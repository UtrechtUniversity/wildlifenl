package models

type ResponseRecord struct {
	ID            string  `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this response."`
	Text          *string `json:"text,omitempty" doc:"The free format text that the user added to this response."`
	QuestionID    string  `json:"questionID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the question this response is for."`
	InteractionID string  `json:"interactionID,omitempty" format:"uuid" writeOnly:"true" required:"true" doc:"The ID of the interaction this response belongs to."`
	AnswerID      *string `json:"answerID,omitempty" format:"uuid" writeOnly:"true" doc:"The ID of the answer that was selected for this response."`
}

type Response struct {
	ResponseRecord
	Question    Question    `json:"question" doc:"The question this response is for."`
	Interaction Interaction `json:"interaction" doc:"The interaction this response belongs to."`
	Answer      *Answer     `json:"answer" doc:"The answer that was selected for this response."`
	Conveyance  *Conveyance `json:"conveyance,omitempty" doc:"The optional conveyance of a message that should be displayed to the user."`
}
