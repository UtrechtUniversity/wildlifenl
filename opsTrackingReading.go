package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewTrackingReadingInput struct {
	Input
	Body *models.TrackingReadingRecord `json:"trackingReading"`
}

type TrackingReadingHolder struct {
	Body *models.TrackingReading `json:"trackingReading"`
}

type TrackingReadingsHolder struct {
	Body []models.TrackingReading `json:"trackingReadings"`
}

type trackingReadingOperations Operations

func newTrackingReadingOperations() *trackingReadingOperations {
	o := trackingReadingOperations{
		Endpoint: "tracking-reading",
	}
	return &o
}

func (o *trackingReadingOperations) RegisterAdd(api huma.API) {
	name := "Add TrackingReading"
	description := "Submit a new tracking reading."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewTrackingReadingInput) (*TrackingReadingHolder, error) {
		trackingReading, err := stores.NewTrackingReadingStore(relationalDB, timeseriesDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}

		// Add Tracking-Reading -> Create Conveyance.
		conveyance, err := stores.NewConveyanceStore(relationalDB).AddForTrackingReading(trackingReading)
		if err != nil {
			return nil, handleError(err)
		}
		trackingReading.Conveyance = conveyance

		return &TrackingReadingHolder{Body: trackingReading}, nil
	})
}
