package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type ConveyanceHolder struct {
	Body *models.Conveyance `json:"conveyance"`
}

type ConveyancesHolder struct {
	Body []models.Conveyance `json:"conveyances"`
}

type conveyanceOperations Operations

func newConveyanceOperations() *conveyanceOperations {
	return &conveyanceOperations{Endpoint: "conveyance"}
}

func (o *conveyanceOperations) RegisterGetAll(api huma.API) {
	name := "Get All Conveyances [deprecated]"
	description := "Retrieve all conveyances. DEPRECATED"
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*ConveyancesHolder, error) {
		conveyances, err := stores.NewConveyanceStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &ConveyancesHolder{Body: conveyances}, nil
	})
}

func (o *conveyanceOperations) RegisterGetMine(api huma.API) {
	name := "Get My Conveyances"
	description := "Retrieve my conveyances."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*ConveyancesHolder, error) {
		conveyances, err := stores.NewConveyanceStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ConveyancesHolder{Body: conveyances}, nil
	})
}

func (o *conveyanceOperations) RegisterGetByExperiment(api huma.API) {
	name := "Get Conveyances By Experiment"
	description := "Retrieve all conveyances for a specific experiment."
	path := "/" + o.Endpoint + "s/experiment/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the experiment to retrieve conveyances for."`
	}) (*ConveyancesHolder, error) {
		conveyances, err := stores.NewConveyanceStore(relationalDB).GetByExperiment(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ConveyancesHolder{Body: conveyances}, nil
	})
}
