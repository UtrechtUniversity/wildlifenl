package wildlifenl

import (
	"context"
	"log"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewInteractionInput struct {
	Input
	Body *models.InteractionRecord `json:"interaction"`
}

type InteractionHolder struct {
	Body *models.Interaction `json:"interaction"`
}

type InteractionsHolder struct {
	Body []models.Interaction `json:"interactions"`
}

type interactionOperations struct{}

func (s *interactionOperations) RegisterGet(api huma.API) {
	scopes := []string{"wildlife-manager", "researcher"}
	huma.Register(api, huma.Operation{
		OperationID: "get-interaction-by-id",
		Tags:        []string{"interaction"},
		Method:      http.MethodGet,
		Path:        "/interaction/{id}",
		Summary:     "Get Interaction By ID",
		Security:    []map[string][]string{{"auth": scopes}},
		Description: "Retrieve a specific interaction by its ID. <br/><br/>**Scopes**<br/>" + scopesAsMarkdown(scopes),
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the interaction." format:"uuid"`
	}) (*InteractionHolder, error) {
		interaction, err := stores.NewInteractionStore(database).Get(input.ID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		if interaction == nil {
			return nil, huma.Error404NotFound("No interaction with ID " + input.ID + " was found")
		}
		return &InteractionHolder{Body: interaction}, nil
	})
}

func (s *interactionOperations) RegisterGetAll(api huma.API) {
	scopes := []string{"researcher"}
	huma.Register(api, huma.Operation{
		OperationID: "get-all-interactions",
		Tags:        []string{"interaction"},
		Method:      http.MethodGet,
		Path:        "/interaction/",
		Summary:     "Get all Interactions",
		Security:    []map[string][]string{{"auth": scopes}},
		Description: "Retrieve all interactions. <br/><br/>**Scopes**<br/>" + scopesAsMarkdown(scopes),
	}, func(ctx context.Context, input *struct{}) (*InteractionsHolder, error) {
		interactions, err := stores.NewInteractionStore(database).GetAll()
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &InteractionsHolder{Body: interactions}, nil
	})
}

func (s *interactionOperations) RegisterAdd(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "add-interaction",
		Tags:        []string{"interaction"},
		Method:      http.MethodPost,
		Path:        "/interaction/",
		Summary:     "Add a new Interaction",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Submit a new Interaction.",
	}, func(ctx context.Context, input *NewInteractionInput) (*InteractionHolder, error) {
		interaction, err := stores.NewInteractionStore(database).Add(input.credential.UserID, input.Body)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &InteractionHolder{Body: interaction}, nil
	})
}
