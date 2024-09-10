package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewAnimalHolder struct {
	Body *models.AnimalRecord `json:"animal"`
}

type AnimalHolder struct {
	Body *models.Animal `json:"animal"`
}

type AnimalsHolder struct {
	Body []models.Animal `json:"animals"`
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
	scopes := []string{"herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this animal." format:"uuid"`
	}) (*AnimalHolder, error) {
		animal, err := stores.NewAnimalStore(relationalDB, timeseriesDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &AnimalHolder{Body: animal}, nil
	})
}

func (o *animalOperations) RegisterGetAll(api huma.API) {
	name := "Get All Animals"
	description := "Retrieve all animals."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"herd-manager", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*AnimalsHolder, error) {
		interactions, err := stores.NewAnimalStore(relationalDB, timeseriesDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &AnimalsHolder{Body: interactions}, nil
	})
}

func (o *animalOperations) RegisterAdd(api huma.API) {
	name := "Add Animal"
	description := "Add a new animal."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"herd-manager"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewAnimalHolder) (*AnimalHolder, error) {
		species, err := stores.NewAnimalStore(relationalDB, timeseriesDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &AnimalHolder{Body: species}, nil
	})
}
