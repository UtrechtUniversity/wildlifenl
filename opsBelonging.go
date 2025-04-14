package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type BelongingHolder struct {
	Body *models.Belonging `json:"belonging"`
}

type BelongingsHolder struct {
	Body []models.Belonging `json:"belonging"`
}

type BelongingUpdateInput struct {
	Input
	ID   string            `path:"id" format:"uuid" doc:"The ID of the belonging to be updated."`
	Body *models.Belonging `json:"belonging"`
}

type belongingOperations Operations

func newBelongingOperations() *belongingOperations {
	return &belongingOperations{Endpoint: "belonging"}
}

func (o *belongingOperations) RegisterGet(api huma.API) {
	name := "Get Belonging By ID"
	description := "Retrieve a specific belonging by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the belonging." format:"uuid"`
	}) (*BelongingHolder, error) {
		belonging, err := stores.NewBelongingStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if belonging == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &BelongingHolder{Body: belonging}, nil
	})
}

func (o *belongingOperations) RegisterGetAll(api huma.API) {
	name := "Get All Belongings"
	description := "Retrieve all belongings."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*BelongingsHolder, error) {
		belongings, err := stores.NewBelongingStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &BelongingsHolder{Body: belongings}, nil
	})
}

func (o *belongingOperations) RegisterAdd(api huma.API) {
	name := "Add Belonging"
	description := "Add a new belonging."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *BelongingHolder) (*BelongingHolder, error) {
		belonging, err := stores.NewBelongingStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &BelongingHolder{Body: belonging}, nil
	})
}

func (o *belongingOperations) RegisterUpdate(api huma.API) {
	name := "Update Belonging"
	description := "Update an existing belonging."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"administrator"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *BelongingUpdateInput) (*BelongingHolder, error) {
		belonging, err := stores.NewBelongingStore(relationalDB).Update(input.ID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &BelongingHolder{Body: belonging}, nil
	})
}
