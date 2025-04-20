package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type TrackingReadingHolder struct {
	Body *models.TrackingReading `json:"trackingReading"`
}

type TrackingReadingsHolder struct {
	Body []models.TrackingReading `json:"trackingReadings"`
}

type TrackingReadingAddInput struct {
	Input
	Body *models.TrackingReadingRecord `json:"trackingReading"`
}

type trackingReadingOperations Operations

func newTrackingReadingOperations() *trackingReadingOperations {
	return &trackingReadingOperations{Endpoint: "tracking-reading"}
}

func (o *trackingReadingOperations) RegisterAdd(api huma.API) {
	name := "Add TrackingReading"
	description := "Submit a new tracking reading."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *TrackingReadingAddInput) (*TrackingReadingHolder, error) {
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

func (o *trackingReadingOperations) RegisterGetMine(api huma.API) {
	name := "Get My TrackingReadings"
	description := "Retrieve my tracking readings."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*TrackingReadingsHolder, error) {
		trackingReadings, err := stores.NewTrackingReadingStore(relationalDB, timeseriesDB).GetForUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &TrackingReadingsHolder{Body: trackingReadings}, nil
	})
}
