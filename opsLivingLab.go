package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type LivingLabHolder struct {
	Body *models.LivingLab `json:"livinglab"`
}

type LivingLabsHolder struct {
	Body []models.LivingLab `json:"livinglabs"`
}

type LivingLabAddInput struct {
	Input
	Body *models.LivingLab `json:"livinglab"`
}

type LivingLabUpdateInput struct {
	Input
	ID   string            `path:"id" format:"uuid" doc:"The ID of the living lab to be updated."`
	Body *models.LivingLab `json:"livinglab"`
}

type livingLabOperations Operations

func newLivingLabOperations() *livingLabOperations {
	return &livingLabOperations{Endpoint: "livinglab"}
}

func (o *livingLabOperations) RegisterGet(api huma.API) {
	name := "Get LivingLab By ID"
	description := "Retrieve a specific living lab by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the living lab."`
	}) (*LivingLabHolder, error) {
		livingLab, err := stores.NewLivingLabStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if livingLab == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &LivingLabHolder{Body: livingLab}, nil
	})
}

func (o *livingLabOperations) RegisterGetAll(api huma.API) {
	name := "Get All LivingLabs"
	description := "Retrieve all living labs."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (output *LivingLabsHolder, err error) {
		livingLabs, err := stores.NewLivingLabStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &LivingLabsHolder{Body: livingLabs}, nil
	})
}

func (o *livingLabOperations) RegisterAdd(api huma.API) {
	name := "Add New LivingLab"
	description := "Submit a new living lab."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *LivingLabAddInput) (*LivingLabHolder, error) {
		if len(input.Body.Definition) < 3 {
			return nil, huma.Error400BadRequest("definition must contain 3 or more points")
		}
		livinglab, err := stores.NewLivingLabStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &LivingLabHolder{Body: livinglab}, nil
	})
}

func (o *livingLabOperations) RegisterUpdate(api huma.API) {
	name := "Update LivingLab"
	description := "Update an existing living lab."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"administrator"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *LivingLabUpdateInput) (*LivingLabHolder, error) {
		if len(input.Body.Definition) < 3 {
			return nil, huma.Error400BadRequest("definition must contain 3 or more points")
		}
		livingLab, err := stores.NewLivingLabStore(relationalDB).Update(input.ID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &LivingLabHolder{Body: livingLab}, nil
	})
}
