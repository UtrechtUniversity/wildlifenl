package wildlifenl

import (
	"context"
	"net/http"
	"regexp"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type QuestionHolder struct {
	Body *models.Question `json:"question"`
}

type QuestionsHolder struct {
	Body []models.Question `json:"questions"`
}

type QuestionAddInput struct {
	Input
	Body *models.QuestionRecord `json:"question"`
}

type QuestionDeleteInput struct {
	Input
	ID string `path:"id" format:"uuid" doc:"The ID of the question to be deleted."`
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
		ID string `path:"id" format:"uuid" doc:"The ID of the question to retrieve."`
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
	}, func(ctx context.Context, input *QuestionAddInput) (*QuestionHolder, error) {
		if input.Body.OpenResponseFormat != nil {
			if _, err := regexp.Compile(*input.Body.OpenResponseFormat); err != nil {
				return nil, huma.Error400BadRequest("Field openResponseFormat must either be not present or contain a regular expression")
			}
		}
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

func (o *questionOperations) RegisterDelete(api huma.API) {
	name := "Delete Question"
	description := "Delete a question."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodDelete
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *QuestionDeleteInput) (*struct{}, error) {
		err := stores.NewQuestionStore(relationalDB).Delete(input.ID, input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return nil, nil
	})
}
