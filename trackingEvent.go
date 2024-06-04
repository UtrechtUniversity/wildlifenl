package wildlifenl

import (
	"context"
	"log"
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

type trackingEventOperations struct{}

func (s *trackingEventOperations) RegisterGetMy(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-my-trackingEvents",
		Tags:        []string{"trackingEvent"},
		Method:      http.MethodGet,
		Path:        "/trackingEvent/me/",
		Summary:     "Get my TrackingEvents",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve all trackingEvents for me.",
	}, func(ctx context.Context, input *Input) (*TrackingEventsHolder, error) {
		trackingEvents, err := stores.NewTrackingEventStore(database).GetByUser(input.credential.UserID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &TrackingEventsHolder{Body: trackingEvents}, nil
	})
}

/*
func (s *trackingEventOperations) RegisterGetByUser(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-trackingEvents-by-user",
		Tags:        []string{"trackingEvent"},
		Method:      http.MethodGet,
		Path:        "/trackingEvent/{userID}",
		Summary:     "Get my TrackingEvents",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve all trackingEvents for a specific user.",
	}, func(ctx context.Context, input *MyTrackingEventsInput) (*TrackingEventsHolder, error) {
		trackingEvents, err := stores.NewTrackingEventStore(database).GetByUser(input.credential.UserID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &TrackingEventsHolder{Body: trackingEvents}, nil
	})
}
*/

func (s *trackingEventOperations) RegisterAdd(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "add-trackingEvent",
		Tags:        []string{"trackingEvent"},
		Method:      http.MethodPost,
		Path:        "/trackingEvent/",
		Summary:     "Add a new TrackingEvent",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Submit a new TrackingEvent.",
	}, func(ctx context.Context, input *NewTrackingEventInput) (*TrackingEventHolder, error) {
		trackingEvent, err := stores.NewTrackingEventStore(database).Add(input.credential.UserID, input.Body)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &TrackingEventHolder{Body: trackingEvent}, nil
	})
}
