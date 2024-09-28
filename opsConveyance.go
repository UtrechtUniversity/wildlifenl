package wildlifenl

import (
	"context"
	"database/sql"
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

func newConveyanceOperations(database *sql.DB) *conveyanceOperations {
	o := conveyanceOperations{
		Database: database,
		Endpoint: "conveyance",
	}
	return &o
}

func (o *conveyanceOperations) RegisterGetAll(api huma.API) {
	name := "Get All Conveyances"
	description := "Retrieve all conveyances."
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
