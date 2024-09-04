package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type InteractionTypesHolder struct {
	Body []models.InteractionType `json:"interactionType"`
}

type interactionTypeOperations Operations

func newInteractionTypeOperations(database *sql.DB) *interactionTypeOperations {
	o := interactionTypeOperations{
		Database: database,
		Endpoint: "interactionType",
	}
	return &o
}

func (o *interactionTypeOperations) RegisterGetAll(api huma.API) {
	name := "Get All InteractionTypes"
	description := "Retrieve all interaction types."
	path := "/" + o.Endpoint + "/"
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
