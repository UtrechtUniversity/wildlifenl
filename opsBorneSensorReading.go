package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewBorneSensorReadingInput struct {
	Body *models.BorneSensorReading `json:"borneSensorReading"`
}

type BorneSensorReadingsHolder struct {
	Body []models.BorneSensorReading `json:"borneSensorReadings"`
}

type borneSensorReadingOperations Operations

func newBorneSensorReadingOperations() *borneSensorReadingOperations {
	o := borneSensorReadingOperations{
		Endpoint: "borne-sensor-reading",
	}
	return &o
}

func (o *borneSensorReadingOperations) RegisterGetAll(api huma.API) {
	name := "Get All BorneSensorReadings"
	description := "Retrieve all borne sensor reading of the last hour."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*BorneSensorReadingsHolder, error) {
		borneSensorReadings, err := stores.NewBorneSensorReadingStore(relationalDB, timeseriesDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorReadingsHolder{Body: borneSensorReadings}, nil
	})
}

func (o *borneSensorReadingOperations) RegisterAdd(api huma.API) {
	name := "Add BorneSensorReading"
	description := "Submit a new reading for a borne sensor."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"data-system"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewBorneSensorReadingInput) (*struct{}, error) {
		err := stores.NewBorneSensorReadingStore(relationalDB, timeseriesDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return nil, nil
	})
}
