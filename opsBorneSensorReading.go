package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type BorneSensorReadingsHolder struct {
	Body []models.BorneSensorReading `json:"borneSensorReadings"`
}

type BorneSensorReadingAddInput struct {
	Input
	Body *models.BorneSensorReadingRecord `json:"borneSensorReading"`
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
	description := "Retrieve all borne sensor reading of the last year."
	path := "/" + o.Endpoint + "s/"
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

func (o *borneSensorReadingOperations) RegisterGetAllBySensor(api huma.API) {
	name := "Get BorneSensorReadings By Sensor"
	description := "Retrieve all borne-sensor readings for a specific sensor."
	path := "/" + o.Endpoint + "s/sensor/{id}"
	scopes := []string{"land-user", "nature-area-manager", "wildlife-manager", "herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the sensor to retrieve borne-sensor readings for."`
	}) (*BorneSensorReadingsHolder, error) {
		borneSensorReadings, err := stores.NewBorneSensorReadingStore(relationalDB, timeseriesDB).GetAllBySensorID(input.ID)
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
	}, func(ctx context.Context, input *BorneSensorReadingAddInput) (*struct{}, error) {
		animal, err := stores.NewBorneSensorReadingStore(relationalDB, timeseriesDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}

		// Add Borne-Sensor-Reading -> Create Alarms.
		if animal != nil {
			ids, err := stores.NewAlarmStore(relationalDB).AddAllFromAnimal(animal)
			if err != nil {
				return nil, handleError(err)
			}

			// From created Alarms -> Create Conveyances
			if err := stores.NewConveyanceStore(relationalDB).AddForAlarmIDs(ids); err != nil {
				return nil, handleError(err)
			}
		}

		return nil, nil
	})
}
