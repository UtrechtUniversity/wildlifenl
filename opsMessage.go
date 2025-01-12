package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type MessageHolder struct {
	Body *models.Message `json:"message"`
}

type MessagesHolder struct {
	Body []models.Message `json:"messages"`
}

type MessageAddInput struct {
	Input
	Body *models.MessageRecord `json:"message"`
}

type MessageDeleteInput struct {
	Input
	ID string `path:"id" format:"uuid" doc:"The ID of the message to be deleted."`
}

type messageOperations Operations

func newMessageOperations() *messageOperations {
	return &messageOperations{Endpoint: "message"}
}

func (o *messageOperations) RegisterGet(api huma.API) {
	name := "Get Message By ID"
	description := "Retrieve a specific message by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this message." format:"uuid"`
	}) (*MessageHolder, error) {
		message, err := stores.NewMessageStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if message == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &MessageHolder{Body: message}, nil
	})
}

func (o *messageOperations) RegisterAdd(api huma.API) {
	name := "Add Message"
	description := "Add a new message."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *MessageAddInput) (*MessageHolder, error) {

		switch input.Body.Trigger {
		case "encounter":
			if input.Body.SpeciesID == nil || input.Body.EncounterMeters == nil || input.Body.EncounterMinutes == nil || input.Body.AnswerID != nil {
				return nil, huma.Error400BadRequest("trigger 'encounter' requires fields speciesID, encounterMeters and encounterMinutes to be non-empty, and field answerID to be empty")
			}
		case "answer":
			if input.Body.SpeciesID != nil || input.Body.EncounterMeters != nil || input.Body.EncounterMinutes != nil || input.Body.AnswerID == nil {
				return nil, huma.Error400BadRequest("trigger 'answer' requires field answerID to be non-empty, and fields speciesID, encounterMeters and encounterMinutes to be empty")
			}
		case "alarm":
			if input.Body.SpeciesID == nil || input.Body.EncounterMeters != nil || input.Body.EncounterMinutes != nil || input.Body.AnswerID != nil {
				return nil, huma.Error400BadRequest("trigger 'alarm' requires field speciesID to be non-empty, and fields answerID, encounterMeters and encounterMinutes to be empty")
			}
		}

		experiments, err := stores.NewExperimentStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		var experiment models.Experiment
		for _, e := range experiments {
			if e.ID == input.Body.ExperimentID {
				experiment = e
				break
			}
		}
		if experiment.ID == "" {
			return nil, generateNotFoundForThisUserError("experiment", input.Body.ExperimentID)
		}
		message, err := stores.NewMessageStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &MessageHolder{Body: message}, nil
	})
}

func (o *messageOperations) RegisterDelete(api huma.API) {
	name := "Delete Message"
	description := "Delete a message."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodDelete
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *MessageDeleteInput) (*struct{}, error) {
		err := stores.NewMessageStore(relationalDB).Delete(input.ID, input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return nil, nil
	})
}

func (o *messageOperations) RegisterGetByExperiment(api huma.API) {
	name := "Get Messages By Experiment"
	description := "Retrieve all messages for a specific experiment."
	path := "/" + o.Endpoint + "s/experiment/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the experiment to retrieve messages for." format:"uuid"`
	}) (*MessagesHolder, error) {
		messages, err := stores.NewMessageStore(relationalDB).GetByExperiment(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &MessagesHolder{Body: messages}, nil
	})
}
