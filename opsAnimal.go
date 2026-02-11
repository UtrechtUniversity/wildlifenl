package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type AnimalHolder struct {
	Body *models.Animal `json:"animal"`
}

type AnimalsHolder struct {
	Body []models.Animal `json:"animals"`
}

type AnimalAddInput struct {
	Body *models.AnimalRecord `json:"animal"`
}

type animalOperations Operations

func newAnimalOperations() *animalOperations {
	return &animalOperations{Endpoint: "animal"}
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
		ID string `path:"id" format:"uuid" doc:"The ID of this animal."`
	}) (*AnimalHolder, error) {
		animal, err := stores.NewAnimalStore(relationalDB, timeseriesDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if animal == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &AnimalHolder{Body: animal}, nil
	})
}

func (o *animalOperations) RegisterGetAll(api huma.API) {
	name := "Get All Animals"
	description := "Retrieve all animals."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
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
	}, func(ctx context.Context, input *AnimalAddInput) (*AnimalHolder, error) {
		species, err := stores.NewAnimalStore(relationalDB, timeseriesDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &AnimalHolder{Body: species}, nil
	})
}

func (o *animalOperations) RegisterGetFilter(api huma.API) {
	name := "Get Animals By Filter"
	description := "Retrieve animals within a spatiotemporal span."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"nature-area-manager", "wildlife-manager", "herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *SpatiotemporalInput) (*AnimalsHolder, error) {
		area := models.Circle{Location: models.Point{Latitude: input.Latitude, Longitude: input.Longitude}, Radius: float64(input.Radius)}
		animals, err := stores.NewAnimalStore(relationalDB, timeseriesDB).GetFiltered(&area, &input.Start, &input.End)
		if err != nil {
			return nil, handleError(err)
		}
		return &AnimalsHolder{Body: animals}, nil
	})
}
