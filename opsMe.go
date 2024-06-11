package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type MeInput struct {
	Input
}

type MeHolder struct {
	Body *models.Me `json:"me"`
}

type meOperations Operations

func newMeOperations(database *sql.DB) *meOperations {
	o := meOperations{
		Database: database,
		Endpoint: "me",
	}
	return &o
}

func (o *meOperations) RegisterGet(api huma.API) {
	name := "Get My Profile"
	description := "Retrieve the current user."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *MeInput) (*MeHolder, error) {
		me, err := stores.NewMeStore(database).Get(input.credential.Token)
		if err != nil {
			return nil, handleError(err)
		}
		if me == nil {
			return nil, huma.Error404NotFound("The logged in user was not found.")
		}
		return &MeHolder{Body: me}, nil
	})
}

func (o *meOperations) RegisterPut(api huma.API) {
	name := "Update My Profile"
	description := "Update the current user."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *MeHolder) (*MeHolder, error) {
		return nil, huma.Error501NotImplemented("")
	})
}
