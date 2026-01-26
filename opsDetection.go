package wildlifenl

import (
	"context"
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

type DetectionAddInput struct {
	Input
	Body models.DetectionRecord `json:"detection"`
}

type detectionOperations Operations

func newDetectionOperations() *detectionOperations {
	return &detectionOperations{Endpoint: "detection"}
}

func (o *detectionOperations) RegisterGetAll(api huma.API) {
	name := "Get All Detections"
	description := "Retrieve all detections."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
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
	description := "Add a new detection.<br/><br/><i>Note that all animals in a detection need to be of the same species. If a single sensor event in the external system identified animals of different species, add multiple detections for that single event here.</i>"
	path := "/" + o.Endpoint + "/"
	scopes := []string{"data-system"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *DetectionAddInput) (*DetectionHolder, error) {
		detection, err := stores.NewDetectionStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}

		// Add Detection -> Create Alarms.
		ids, err := stores.NewAlarmStore(relationalDB).AddAllFromDetection(detection)
		if err != nil {
			return nil, handleError(err)
		}

		// From created Alarms -> Create Conveyances
		if err := stores.NewConveyanceStore(relationalDB).AddForAlarmIDs(ids); err != nil {
			return nil, handleError(err)
		}

		return &DetectionHolder{Body: detection}, nil
	})
}
