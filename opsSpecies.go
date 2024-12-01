package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type SpeciesHolder struct {
	Body *models.Species `json:"species"`
}

type SpecieSHolder struct {
	Body []models.Species `json:"species"`
}

type UpdateSpeciesInput struct {
	Input
	ID   string          `query:"ID" format:"uuid" doc:"The ID of the species to be updated."`
	Body *models.Species `json:"species"`
}

type speciesOperations Operations

func newSpeciesOperations() *speciesOperations {
	return &speciesOperations{Endpoint: "species"}
}

func (o *speciesOperations) RegisterGet(api huma.API) {
	name := "Get Species By ID"
	description := "Retrieve a specific species by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the species." format:"uuid"`
	}) (*SpeciesHolder, error) {
		species, err := stores.NewSpeciesStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if species == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &SpeciesHolder{Body: species}, nil
	})
}

func (o *speciesOperations) RegisterGetAll(api huma.API) {
	name := "Get All Species"
	description := "Retrieve all species."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*SpecieSHolder, error) {
		speciesX, err := stores.NewSpeciesStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &SpecieSHolder{Body: speciesX}, nil
	})
}

func (o *speciesOperations) RegisterAdd(api huma.API) {
	name := "Add Species"
	description := "Add a new species."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *SpeciesHolder) (*SpeciesHolder, error) {
		species, err := stores.NewSpeciesStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &SpeciesHolder{Body: species}, nil
	})
}

func (o *speciesOperations) RegisterUpdate(api huma.API) {
	name := "Update Species"
	description := "Update an existing species."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *UpdateSpeciesInput) (*SpeciesHolder, error) {
		species, err := stores.NewSpeciesStore(relationalDB).Update(input.ID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &SpeciesHolder{Body: species}, nil
	})
}
