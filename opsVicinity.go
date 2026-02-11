package wildlifenl

import (
	"context"
	"net/http"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type VicinityHolder struct {
	Body *models.Vicinity `json:"vicinity"`
}

type vicinityOperations Operations

func newVicinityOperations() *vicinityOperations {
	return &vicinityOperations{Endpoint: "vicinity"}
}

func (o *vicinityOperations) RegisterGetMine(api huma.API) {
	name := "Get My Vicinity"
	description := "Retrieve all Animals, Interactions and Detections in the vicinity of the current user. Vicinity is defined as a circular area of ~500 metres diameter around the user's current location, and a timespan of 48 hours since now."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*VicinityHolder, error) {
		const radius int = 250
		const hours int = 48

		vicinity := new(models.Vicinity)
		profile, err := stores.NewProfileStore(relationalDB).Get(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		if profile.Location != nil {
			area := models.Circle{
				Location: models.Point{
					Latitude:  profile.Location.Latitude,
					Longitude: profile.Location.Longitude,
				},
				Radius: float64(radius),
			}
			before := time.Now()
			after := before.Add(-time.Duration(hours) * time.Hour)

			animals, err := stores.NewAnimalStore(relationalDB, timeseriesDB).GetFiltered(&area, &before, &after)
			if err != nil {
				return nil, handleError(err)
			}
			vicinity.Animals = animals

			interactions, err := stores.NewInteractionStore(relationalDB).GetFiltered(&area, &before, &after)
			if err != nil {
				return nil, handleError(err)
			}
			vicinity.Interactions = interactions

			detections, err := stores.NewDetectionStore(relationalDB).GetFiltered(&area, &after, &before)
			if err != nil {
				return nil, handleError(err)
			}
			vicinity.Detections = detections
		}
		return &VicinityHolder{Body: vicinity}, nil
	})
}
