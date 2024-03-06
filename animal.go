package wildlifenl

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
)

type Animal struct {
	ID      int     `json:"id" example:"42" minimum:"1" doc:"The identification number of this animal."`
	Name    string  `json:"name" example:"Flupke" doc:"The name of this animal."`
	Species Species `json:"species"`
}

type AnimalResult struct {
	Body *Animal `json:"animal" doc:"A specific animal."`
}

type animalOperations struct{}

func (s *animalOperations) RegisterAnimalGet(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-animal-by-id",
		Method:      http.MethodGet,
		Path:        "/animal/{id}",
		Summary:     "Animal by ID",
		Security: []map[string][]string{
			{"auth": {"researcher", "herd-manager"}},
		},
		Description: "Retrieve a specific animal by its ID.<br/><br/>**Scopes**<br/>`researcher`, `herd-manager`",
	}, animalGet)
}

func animalGet(ctx context.Context, input *struct {
	ID int `path:"id" example:"42" doc:"The identification number for the animal to retrieve."`
}) (*AnimalResult, error) {
	animal := app.store.getAnimal(input.ID)
	if animal == nil {
		return nil, huma.Error404NotFound("No animal with id " + strconv.Itoa(input.ID) + " was found")
	}
	output := new(AnimalResult)
	output.Body = animal
	return output, nil
}
