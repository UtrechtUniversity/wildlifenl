package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewTrackingEventInput struct {
	Input
	Body *models.TrackingEventRecord `json:"trackingEvent"`
}

type TrackingEventHolder struct {
	Body *models.TrackingEvent `json:"trackingEvent"`
}

type TrackingEventsHolder struct {
	Body []models.TrackingEvent `json:"trackingEvents"`
}

type trackingEventOperations Operations

func newTrackingEventOperations(database *sql.DB) *trackingEventOperations {
	o := trackingEventOperations{
		Database: database,
		Endpoint: "tracking-event",
	}
	return &o
}

func (o *trackingEventOperations) RegisterGetMy(api huma.API) {
	name := "Get my TrackingEvents"
	description := "Retrieve all trackingEvents for me."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*TrackingEventsHolder, error) {
		trackingEvents, err := stores.NewTrackingEventStore(o.Database).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &TrackingEventsHolder{Body: trackingEvents}, nil
	})
}

func (o *trackingEventOperations) RegisterAdd(api huma.API) {
	name := "Add a new TrackingEvent"
	description := "Submit a new TrackingEvent."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewTrackingEventInput) (*TrackingEventHolder, error) {
		trackingEvent, err := stores.NewTrackingEventStore(o.Database).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &TrackingEventHolder{Body: trackingEvent}, nil
	})
}
