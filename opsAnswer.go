package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type AnswerHolder struct {
	Body *models.Answer `json:"answer"`
}

type AnswersHolder struct {
	Body []models.Answer `json:"answers"`
}

type NewAnswerInput struct {
	Input
	Body *models.AnswerRecord `json:"answer"`
}

type AnswerDeleteInput struct {
	Input
	ID string `path:"id" format:"uuid" doc:"The ID of the answer to be deleted."`
}

type answerOperations Operations

func newAnswerOperations() *answerOperations {
	return &answerOperations{Endpoint: "answer"}
}

func (o *answerOperations) RegisterGet(api huma.API) {
	name := "Get Answer By ID"
	description := "Retrieve a specific answer by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this answer." format:"uuid"`
	}) (*AnswerHolder, error) {
		answer, err := stores.NewAnswerStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if answer == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &AnswerHolder{Body: answer}, nil
	})
}

func (o *answerOperations) RegisterAdd(api huma.API) {
	name := "Add Answer"
	description := "Add a new answer."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewAnswerInput) (*AnswerHolder, error) {
		if input.Body.NextQuestionID != nil && *input.Body.NextQuestionID == input.Body.QuestionID {
			return nil, huma.Error400BadRequest("Fields questionID and nextQuestionID must have different values")
		}
		store := stores.NewQuestionnaireStore(relationalDB)
		questionnaires, err := store.GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		var question models.Question
		var nextQuestion models.Question
		for _, r := range questionnaires {
			for _, q := range r.Questions {
				if q.ID == input.Body.QuestionID {
					question = q
				}
				if input.Body.NextQuestionID != nil && q.ID == *input.Body.NextQuestionID {
					nextQuestion = q
				}
			}
		}
		if question.ID == "" {
			return nil, generateNotFoundForThisUserError("question", input.Body.QuestionID)
		}
		if input.Body.NextQuestionID != nil {
			if nextQuestion.ID == "" {
				return nil, generateNotFoundForThisUserError("question", *input.Body.NextQuestionID)
			}
			if nextQuestion.QuestionnaireID != question.QuestionnaireID {
				return nil, huma.Error400BadRequest("Fields questionID and nextQuestionID must both refer to questions in the same questionnaire.")
			}
		}
		answer, err := stores.NewAnswerStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &AnswerHolder{Body: answer}, nil
	})
}

func (o *answerOperations) RegisterDelete(api huma.API) {
	name := "Delete Answer"
	description := "Delete an anwser."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodDelete
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *AnswerDeleteInput) (*struct{}, error) {
		err := stores.NewAnswerStore(relationalDB).Delete(input.ID, input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return nil, nil
	})
}
