package wildlifenl

import (
	"context"
	"net/http"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type InteractionHolder struct {
	Body *models.Interaction `json:"interaction"`
}

type InteractionsHolder struct {
	Body []models.Interaction `json:"interactions"`
}

type InteractionQueryInput struct {
	Latitude  float64   `query:"area_latitude" minimum:"-90" maximum:"90" doc:"The latitude of the point that defines the centre of the area in which the interaction happened."`
	Longitude float64   `query:"area_longitude" minimum:"-180" maximum:"180" doc:"The longitude of the point that defines the centre of the area in which the interaction happened."`
	Radius    int       `query:"area_radius" minimum:"1" doc:"The radius of the area as defined by its centre point, in which the interaction happened."`
	Before    time.Time `query:"moment_before" doc:"The moment in time before which the interaction happened."`
	After     time.Time `query:"moment_after" doc:"The moment in time after which the interaction happened."`
}

type InteractionAddInput struct {
	Input
	Body *models.InteractionRecord `json:"interaction"`
}

type interactionOperations Operations

func newInteractionOperations() *interactionOperations {
	return &interactionOperations{Endpoint: "interaction"}
}

func (o *interactionOperations) RegisterGet(api huma.API) {
	name := "Get Interaction By ID"
	description := "Retrieve a specific interaction by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"wildlife-manager", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the interaction." format:"uuid"`
	}) (*InteractionHolder, error) {
		interaction, err := stores.NewInteractionStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if interaction == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &InteractionHolder{Body: interaction}, nil
	})
}

func (o *interactionOperations) RegisterGetAll(api huma.API) {
	name := "Get All Interactions"
	description := "Retrieve all interactions."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*InteractionsHolder, error) {
		interactions, err := stores.NewInteractionStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}

func (o *interactionOperations) RegisterGetMine(api huma.API) {
	name := "Get My Interactions"
	description := "Retrieve my interactions."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*InteractionsHolder, error) {
		interactions, err := stores.NewInteractionStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}

func (o *interactionOperations) RegisterAdd(api huma.API) {
	name := "Add Interaction"
	description := "Submit a new interaction."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *InteractionAddInput) (*InteractionHolder, error) {
		if input.Body.TypeID == 1 && input.Body.ReportOfSighting == nil {
			return nil, huma.Error400BadRequest("Interaction of TypeID=1 must contain a report of sighting")
		}
		if input.Body.TypeID == 2 && input.Body.ReportOfDamage == nil {
			return nil, huma.Error400BadRequest("Interaction of TypeID=2 must contain a report of damage")
		}
		if input.Body.TypeID == 3 && input.Body.ReportOfCollision == nil {
			return nil, huma.Error400BadRequest("Interaction of TypeID=3 must contain a report of collision")
		}

		interaction, err := stores.NewInteractionStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}

		// Add Interaction -> Get Questionnaire.
		questionnaire, err := stores.NewQuestionnaireStore(relationalDB).AssignRandomToInteraction(interaction)
		if err != nil {
			return nil, handleError(err)
		}
		if questionnaire != nil {
			interaction.Questionnaire = questionnaire
		}

		// Add Interaction -> Create Alarms.
		if interaction.Type.ID == 1 {
			ids, err := stores.NewAlarmStore(relationalDB).AddAllFromInteraction(interaction)
			if err != nil {
				return nil, handleError(err)
			}

			// From created Alarms -> Create Conveyances
			if err := stores.NewConveyanceStore(relationalDB).AddForAlarmIDs(ids); err != nil {
				return nil, handleError(err)
			}
		}

		return &InteractionHolder{Body: interaction}, nil
	})
}

func (o *interactionOperations) RegisterGetFilter(api huma.API) {
	name := "Get Interactions By Filter"
	description := "Retrieve interactions within a spatiotemporal span."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"nature-area-manager", "wildlife-manager", "herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *SpatiotemporalInput) (*InteractionsHolder, error) {
		area := models.Circle{Location: models.Point{Latitude: input.Latitude, Longitude: input.Longitude}, Radius: float64(input.Radius)}
		interactions, err := stores.NewInteractionStore(relationalDB).GetFiltered(&area, &input.Start, &input.End)
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}
