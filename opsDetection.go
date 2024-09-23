package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type DetectionHolder struct {
	Body *models.Detection `json:"detection"`
}

type DetectionsHolder struct {
	Body []models.Detection `json:"detections"`
}

type NewDetectionInput struct {
	Body models.DetectionRecord `json:"detection"`
}

type detectionOperations Operations

func newDetectionOperations(database *sql.DB) *detectionOperations {
	o := detectionOperations{
		Database: database,
		Endpoint: "detection",
	}
	return &o
}

func (o *detectionOperations) RegisterGetAll(api huma.API) {
	name := "Get All Detections"
	description := "Retrieve all detections."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*DetectionsHolder, error) {
		detections, err := stores.NewDetectionStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &DetectionsHolder{Body: detections}, nil
	})
}

func (o *detectionOperations) RegisterAdd(api huma.API) {
	name := "Add Detection"
	description := "Add a new detection."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"data-system"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewDetectionInput) (*DetectionHolder, error) {
		detection, err := stores.NewDetectionStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		if err := stores.NewAlarmStore(relationalDB).AddAllFromDetection(detection); err != nil {
			return nil, handleError(err)
		}
		return &DetectionHolder{Body: detection}, nil
	})
}
