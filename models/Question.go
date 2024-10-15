package models

type QuestionRecord struct {
	ID                    string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this question."`
	Text                  string `json:"text" doc:"The text of this question."`
	Description           string `json:"description" doc:"The further explanation for this questions."`
	Index                 int    `json:"index" doc:"The index of this question within the questionnaire. If multiple questions have the same index, the client application should present them to the end user in a randomly shuffled order."`
	AllowMultipleResponse bool   `json:"allowMultipleResponse" doc:"Whether or not this questions allows for multiple reponses."`
	AllowOpenResponse     bool   `json:"allowOpenResponse" doc:"Whether or not this question allows for a response to contain a free format text."`
	QuestionnaireID       string `json:"questionnaireID,omitempty" format:"uuid" writeOnly:"true" required:"true" docs:"The ID of the questionnaire this question belongs to."`
}

type Question struct {
	QuestionRecord
	Answers []Answer `json:"answers,omitempty" doc:"The selectable answers for this questions."`
}
