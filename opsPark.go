package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewParkInput struct {
	Input
	Body *models.Park `json:"park"`
}

type ParkHolder struct {
	Body *models.Park `json:"park"`
}

type ParksHolder struct {
	Body []models.Park `json:"parks"`
}

type parkOperations Operations

func newParkOperations(database *sql.DB) *parkOperations {
	o := parkOperations{
		Database: database,
		Endpoint: "park",
	}
	return &o
}

func (o *parkOperations) RegisterGet(api huma.API) {
	name := "Get Park By ID"
	description := "Retrieve a specific park by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the park."`
	}) (*ParkHolder, error) {
		park, err := stores.NewParkStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if park == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &ParkHolder{Body: park}, nil
	})
}

func (o *parkOperations) RegisterGetAll(api huma.API) {
	name := "Get all Parks"
	description := "Retrieve all parks."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (output *ParksHolder, err error) {
		parks, err := stores.NewParkStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &ParksHolder{Body: parks}, nil
	})
}

func (o *parkOperations) RegisterAdd(api huma.API) {
	name := "Add New Park"
	description := "Submit a new park."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewParkInput) (*ParkHolder, error) {
		park, err := stores.NewParkStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &ParkHolder{Body: park}, nil
	})
}
