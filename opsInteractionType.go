package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type InteractionTypeHolder struct {
	Body *models.InteractionType `json:"interactionType"`
}

type InteractionTypesHolder struct {
	Body []models.InteractionType `json:"interactionTypes"`
}

type InteractionTypeUpdateInput struct {
	ID   int                     `path:"id" doc:"The ID of the interaction type to be updated."`
	Body *models.InteractionType `json:"interactionType"`
}

type interactionTypeOperations Operations

func newInteractionTypeOperations() *interactionTypeOperations {
	return &interactionTypeOperations{Endpoint: "interactionType"}
}

func (o *interactionTypeOperations) RegisterGetAll(api huma.API) {
	name := "Get All InteractionTypes"
	description := "Retrieve all interaction types."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*InteractionTypesHolder, error) {
		interactionTypes, err := stores.NewInteractionTypeStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionTypesHolder{Body: interactionTypes}, nil
	})
}

func (o *interactionTypeOperations) RegisterAdd(api huma.API) {
	name := "Add Interaction Type"
	description := "Add a new interaction type."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *InteractionTypeHolder) (*InteractionTypeHolder, error) {
		interactionType, err := stores.NewInteractionTypeStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionTypeHolder{Body: interactionType}, nil
	})
}

func (o *interactionTypeOperations) RegisterUpdate(api huma.API) {
	name := "Update Interaction Type"
	description := "Update an existing interaction type."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"administrator"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *InteractionTypeUpdateInput) (*InteractionTypeHolder, error) {
		interactionType, err := stores.NewInteractionTypeStore(relationalDB).Update(input.ID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionTypeHolder{Body: interactionType}, nil
	})
}
