package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewQuestionInput struct {
	Input
	Body *models.QuestionRecord `json:"question"`
}

type QuestionHolder struct {
	Body *models.Question `json:"question"`
}

type QuestionsHolder struct {
	Body []models.Question `json:"questions"`
}

type questionOperations Operations

func newQuestionOperations() *questionOperations {
	return &questionOperations{Endpoint: "question"}
}

func (o *questionOperations) RegisterGet(api huma.API) {
	name := "Get Question By ID"
	description := "Retrieve a specific question by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this question." format:"uuid"`
	}) (*QuestionHolder, error) {
		question, err := stores.NewQuestionStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if question == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &QuestionHolder{Body: question}, nil
	})
}

func (o *questionOperations) RegisterAdd(api huma.API) {
	name := "Add Question"
	description := "Add a new question."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewQuestionInput) (*QuestionHolder, error) {
		store := stores.NewQuestionnaireStore(relationalDB)
		questionnaires, err := store.GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		var questionnaire models.Questionnaire
		for _, q := range questionnaires {
			if q.ID == input.Body.QuestionnaireID {
				questionnaire = q
				break
			}
		}
		if questionnaire.ID == "" {
			return nil, generateNotFoundForThisUserError("questionnaire", input.Body.QuestionnaireID)
		}
		question, err := stores.NewQuestionStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &QuestionHolder{Body: question}, nil
	})
}
