package wildlifenl

import (
	"context"
	"net/http"
	"regexp"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type ResponseHolder struct {
	Body *models.Response `json:"response"`
}

type ResponsesHolder struct {
	Body []models.Response `json:"responses"`
}

type ResponseAddInput struct {
	Input
	Body *models.ResponseRecord `json:"response"`
}

type responseOperations Operations

func newResponseOperations() *responseOperations {
	return &responseOperations{Endpoint: "response"}
}

func (o *responseOperations) RegisterGetByQuestionnaire(api huma.API) {
	name := "Get Responses by Questionnaire"
	description := "Retrieve responses for a specific questionnaire."
	path := "/" + o.Endpoint + "s/questionnaire/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the questionnaire to retrieve responses for."`
	}) (*ResponsesHolder, error) {
		responses, err := stores.NewResponseStore(relationalDB).GetByQuestionnaire(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ResponsesHolder{Body: responses}, nil
	})
}

func (o *responseOperations) RegisterAdd(api huma.API) {
	name := "Add Response"
	description := "Add a new response."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *ResponseAddInput) (*ResponseHolder, error) {

		question, err := stores.NewQuestionStore(relationalDB).Get(input.Body.QuestionID)
		if err != nil {
			return nil, handleError(err)
		}
		if !question.AllowOpenResponse && input.Body.Text != nil {
			return nil, huma.Error400BadRequest("question (" + question.ID + ") does not allow open responses, therefore field 'text' must not be present")
		}
		if question.AllowOpenResponse && question.OpenResponseFormat != nil {
			r, err := regexp.Compile(*question.OpenResponseFormat)
			if err != nil {
				return nil, handleError(err)
			}
			text := *input.Body.Text
			if loc := r.FindStringIndex(text); loc == nil || text[loc[0]:loc[1]] != text {
				return nil, huma.Error400BadRequest("Field text does not match regular expression " + *question.OpenResponseFormat)
			}
		}
		earlierResponses, err := stores.NewResponseStore(relationalDB).GetForInteractionByQuestion(input.Body.InteractionID, input.Body.QuestionID)
		if err != nil {
			return nil, handleError(err)
		}
		if !question.AllowMultipleResponse && len(earlierResponses) > 0 {
			return nil, huma.Error400BadRequest("question (" + question.ID + ") does not allow multiple responses, and a previous response already exists")
		}

		response, err := stores.NewResponseStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		if response == nil {
			return nil, huma.Error400BadRequest("a response could not be added for the combination of interactionID, questionID and answerID that was provided")
		}

		// Add Response -> Create Conveyance.
		conveyance, err := stores.NewConveyanceStore(relationalDB).AddForResponse(response)
		if err != nil {
			return nil, handleError(err)
		}
		response.Conveyance = conveyance

		return &ResponseHolder{Body: response}, nil
	})
}

func (o *responseOperations) RegisterGetByExperiment(api huma.API) {
	name := "Get Responses by Experiment"
	description := "Retrieve responses for a specific experiment."
	path := "/" + o.Endpoint + "s/experiment/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the experiment to retrieve responses for."`
	}) (*ResponsesHolder, error) {
		responses, err := stores.NewResponseStore(relationalDB).GetByExperiment(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ResponsesHolder{Body: responses}, nil
	})
}
