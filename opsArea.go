package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewAreaInput struct {
	Input
	Body *models.AreaRecord `json:"area"`
}

type AreaHolder struct {
	Body *models.Area `json:"area"`
}

type AreasHolder struct {
	Body []models.Area `json:"area"`
}

type areaOperations Operations

func newAreaOperations(database *sql.DB) *areaOperations {
	o := areaOperations{
		Database: database,
		Endpoint: "area",
	}
	return &o
}

func (o *areaOperations) RegisterGet(api huma.API) {
	name := "Get Area By ID"
	description := "Retrieve a specific area by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the area."`
	}) (*AreaHolder, error) {
		area, err := stores.NewAreaStore(database).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if area == nil {
			return nil, huma.Error404NotFound("No area with ID " + input.ID + " was found")
		}
		return &AreaHolder{Body: area}, nil
	})
}

func (o *areaOperations) RegisterGetMy(api huma.API) {
	name := "Get My Areas"
	description := "Retrieve all areas made by the current user."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*AreasHolder, error) {
		areas, err := stores.NewAreaStore(database).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &AreasHolder{Body: areas}, nil
	})
}

func (o *areaOperations) RegisterAdd(api huma.API) {
	name := "Add New Area"
	description := "Submit a new area."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewAreaInput) (*AreaHolder, error) {
		area, err := stores.NewAreaStore(database).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &AreaHolder{Body: area}, nil
	})
}
