package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

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

func (o *parkOperations) RegisterGetAll(api huma.API) {
	name := "Get all Parks"
	description := "Retrieve all parks."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (output *ParksHolder, err error) {
		parks, err := stores.NewParkStore(database).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &ParksHolder{Body: parks}, nil
	})
}
