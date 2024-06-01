package wildlifenl

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
)

type AnimalHolder struct {
	Body *models.Animal `json:"animal" doc:"A specific animal."`
}

type animalOperations struct{}

func (s *animalOperations) RegisterGet(api huma.API) {
	scopes := []string{"wildlife-manager", "researcher"}
	huma.Register(api, huma.Operation{
		OperationID: "get-animal-by-id",
		Tags:        []string{"animal"},
		Method:      http.MethodGet,
		Path:        "/animal/{id}",
		Summary:     "Get Animal by ID",
		Security:    []map[string][]string{{"auth": scopes}},
		Description: "Retrieve a specific animal by its ID. <br/><br/>**Scopes**<br/>" + scopesAsMarkdown(scopes),
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this animal." format:"uuid"`
	}) (*AnimalHolder, error) {
		animal, err := stores.NewAnimals(database).Get(input.ID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		if animal == nil {
			return nil, huma.Error404NotFound("No animal with ID " + input.ID + " was found")
		}
		return &AnimalHolder{Body: animal}, nil
	})
}
