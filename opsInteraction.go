package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewInteractionInput struct {
	Input
	Body *models.InteractionRecord `json:"interaction"`
}

type InteractionHolder struct {
	Body *models.Interaction `json:"interaction"`
}

type InteractionsHolder struct {
	Body []models.Interaction `json:"interactions"`
}

type interactionOperations Operations

func newInteractionOperations(database *sql.DB) *interactionOperations {
	o := interactionOperations{
		Database: database,
		Endpoint: "interaction",
	}
	return &o
}

func (o *interactionOperations) RegisterGet(api huma.API) {
	name := "Get Interaction By ID"
	description := "Retrieve a specific interaction by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"wildlife-manager", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the interaction." format:"uuid"`
	}) (*InteractionHolder, error) {
		interaction, err := stores.NewInteractionStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if interaction == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &InteractionHolder{Body: interaction}, nil
	})
}

func (o *interactionOperations) RegisterGetAll(api huma.API) {
	name := "Get All Interactions"
	description := "Retrieve all interactions."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*InteractionsHolder, error) {
		interactions, err := stores.NewInteractionStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}

func (o *interactionOperations) RegisterAdd(api huma.API) {
	name := "Add Interaction"
	description := "Submit a new interaction."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewInteractionInput) (*InteractionHolder, error) {
		interaction, err := stores.NewInteractionStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		questionnaire, err := stores.NewQuestionnaireStore(relationalDB).GetRandomActiveByInteractionType(interaction.Type)
		if err != nil {
			return nil, handleError(err)
		}
		interaction.Questionnaire = questionnaire
		if interaction.Type.ID == 1 { // TODO issue #15: Uses magic number.
			if err := stores.NewAlarmStore(relationalDB).AddAllFromInteraction(interaction); err != nil {
				return nil, handleError(err)
			}
		}
		return &InteractionHolder{Body: interaction}, nil
	})
}

func (o *interactionOperations) RegisterGetMine(api huma.API) {
	name := "Get My Interactions"
	description := "Retrieve my interactions."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*InteractionsHolder, error) {
		interactions, err := stores.NewInteractionStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}
