package wildlifenl

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
)

type SpeciesHolder struct {
	Body *models.Species `json:"species"`
}

type SpeciesXHolder struct {
	Body []models.Species `json:"species"`
}

type speciesOperations struct{}

func (s *speciesOperations) RegisterGet(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-species-by-id",
		Tags:        []string{"species"},
		Method:      http.MethodGet,
		Path:        "/species/{id}",
		Summary:     "Get Species By ID",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve a specific species by ID.",
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the species." format:"uuid"`
	}) (*SpeciesHolder, error) {
		species, err := stores.NewSpeciesStore(database).Get(input.ID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		if species == nil {
			return nil, huma.Error404NotFound("No species with ID " + input.ID + " was found")
		}
		return &SpeciesHolder{Body: species}, nil
	})
}

func (s *speciesOperations) RegisterGetAll(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-all-species",
		Tags:        []string{"species"},
		Method:      http.MethodGet,
		Path:        "/species/",
		Summary:     "Get all Species",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve all species.",
	}, func(ctx context.Context, input *struct{}) (*SpeciesXHolder, error) {
		speciesX, err := stores.NewSpeciesStore(database).GetAll()
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &SpeciesXHolder{Body: speciesX}, nil
	})
}
