package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewMessageInput struct {
	Input
	Body *models.MessageRecord `json:"message"`
}

type MessageHolder struct {
	Body *models.Message `json:"message"`
}

type MessagesHolder struct {
	Body []models.Message `json:"messages"`
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

func (o *messageOperations) RegisterGetAll(api huma.API) {
	name := "Get All Messages"
	description := "Retrieve all messages."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"administrator"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*MessagesHolder, error) {
		messages, err := stores.NewMessageStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &MessagesHolder{Body: messages}, nil
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
	}, func(ctx context.Context, input *NewMessageInput) (*MessageHolder, error) {

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

func (o *messageOperations) RegisterGetMine(api huma.API) {
	name := "Get My Messages"
	description := "Retrieve my messages."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*MessagesHolder, error) {
		messages, err := stores.NewMessageStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &MessagesHolder{Body: messages}, nil
	})
}
