package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type AnimalHolder struct {
	Body *models.Animal `json:"animal" doc:"A specific animal."`
}

type animalOperations Operations

func newAnimalOperations(database *sql.DB) *animalOperations {
	o := animalOperations{
		Database: database,
		Endpoint: "animal",
	}
	return &o
}

func (o *animalOperations) RegisterGet(api huma.API) {
	name := "Get Animal By ID"
	description := "Retrieve a specific animal by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"wildlife-manager", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this animal." format:"uuid"`
	}) (*AnimalHolder, error) {
		animal, err := stores.NewAnimals(database).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &AnimalHolder{Body: animal}, nil
	})
}
