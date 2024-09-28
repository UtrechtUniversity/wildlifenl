package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type EncounterHolder struct {
	Body *models.Encounter `json:"encounter"`
}

type EncountersHolder struct {
	Body []models.Encounter `json:"encounters"`
}

type encounterOperations Operations

func newEncounterOperations(database *sql.DB) *encounterOperations {
	o := encounterOperations{
		Database: database,
		Endpoint: "encounter",
	}
	return &o
}

func (o *encounterOperations) RegisterGetAll(api huma.API) {
	name := "Get All Encounters"
	description := "Retrieve all encounters."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*EncountersHolder, error) {
		encounters, err := stores.NewEncounterStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &EncountersHolder{Body: encounters}, nil
	})
}
