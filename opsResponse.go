package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewResponseInput struct {
	Input
	Body *models.ResponseRecord `json:"response"`
}

type ResponseHolder struct {
	Body *models.Response `json:"response"`
}

type ResponsesHolder struct {
	Body []models.Response `json:"responses"`
}

type responseOperations Operations

func newResponseOperations(database *sql.DB) *responseOperations {
	o := responseOperations{
		Database: database,
		Endpoint: "response",
	}
	return &o
}

func (o *responseOperations) RegisterGet(api huma.API) {
	name := "Get Response By ID"
	description := "Retrieve a specific response by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"administrator"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this response." format:"uuid"`
	}) (*ResponseHolder, error) {
		response, err := stores.NewResponseStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if response == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &ResponseHolder{Body: response}, nil
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
	}, func(ctx context.Context, input *NewResponseInput) (*ResponseHolder, error) {

		// TODO issue 19: do sanity check here.

		// TODO issue 20: check question settings here.

		response, err := stores.NewResponseStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
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
