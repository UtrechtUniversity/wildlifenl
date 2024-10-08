package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewLivingLabInput struct {
	Input
	Body *models.LivingLab `json:"livinglab"`
}

type LivingLabHolder struct {
	Body *models.LivingLab `json:"livinglab"`
}

type LivingLabsHolder struct {
	Body []models.LivingLab `json:"livinglabs"`
}

type livinglabOperations Operations

func newLivingLabOperations() *livinglabOperations {
	return &livinglabOperations{Endpoint: "livinglab"}
}

func (o *livinglabOperations) RegisterGet(api huma.API) {
	name := "Get LivingLab By ID"
	description := "Retrieve a specific living lab by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the livinglab."`
	}) (*LivingLabHolder, error) {
		livinglab, err := stores.NewLivingLabStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if livinglab == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &LivingLabHolder{Body: livinglab}, nil
	})
}

func (o *livinglabOperations) RegisterGetAll(api huma.API) {
	name := "Get all LivingLabs"
	description := "Retrieve all living labs."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (output *LivingLabsHolder, err error) {
		livinglabs, err := stores.NewLivingLabStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &LivingLabsHolder{Body: livinglabs}, nil
	})
}

func (o *livinglabOperations) RegisterAdd(api huma.API) {
	name := "Add New LivingLab"
	description := "Submit a new living lab."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewLivingLabInput) (*LivingLabHolder, error) {
		livinglab, err := stores.NewLivingLabStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &LivingLabHolder{Body: livinglab}, nil
	})
}
