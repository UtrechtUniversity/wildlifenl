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
	Latitude  float64   `query:"area_latitude" minimum:"-90" maximum:"90"`
	Longitude float64   `query:"area_longitude" minimum:"-180" maximum:"180"`
	Radius    int       `query:"area_radius" minimum:"1"`
	Before    time.Time `query:"moment_before"`
	After     time.Time `query:"moment_after"`
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

		// Sanity check to validate the BelongingID of DamageReport before inserting the new Interaction because InteractionStore.Add cannot do this, see InteractionStore.Add().
		if input.Body.TypeID == 2 {
			belonging, err := stores.NewBelongingStore(relationalDB).Get(input.Body.ReportOfDamage.Belonging.ID)
			if err != nil {
				return nil, handleError(err)
			}
			if belonging == nil {
				return nil, generateNotFoundByIDError("Belonging", input.Body.ReportOfDamage.Belonging.ID)
			}
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

func (o *interactionOperations) RegisterQuery(api huma.API) {
	name := "Query Interactions"
	description := "Retrieve interactions for a certain area and/or timespan."
	path := "/" + o.Endpoint + "s/query/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *InteractionQueryInput) (*InteractionsHolder, error) {
		if input.Radius == 0 && input.Before.Year() == 1 && input.After.Year() == 1 {
			return nil, huma.Error400BadRequest("Either all values for `area` or `moment_before` or `moment_after` is required.")
		}
		var area *models.Circle
		if input.Radius > 0 {
			area = &models.Circle{
				Location: models.Point{
					Latitude:  input.Latitude,
					Longitude: input.Longitude,
				},
				Radius: float64(input.Radius),
			}
		}
		var before *time.Time
		if input.Before.Year() > 1 {
			before = &input.Before
		}
		var after *time.Time
		if input.After.Year() > 1 {
			after = &input.After
		}
		interactions, err := stores.NewInteractionStore(relationalDB).GetFiltered(area, before, after)
		if err != nil {
			return nil, handleError(err)
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}
