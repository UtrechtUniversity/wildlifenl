package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type ProfileHolder struct {
	Body *models.Profile `json:"profile"`
}

type ProfilesHolder struct {
	Body []models.Profile `json:"profiles"`
}

type profileOperations Operations

func newProfileOperations() *profileOperations {
	return &profileOperations{Endpoint: "profile"}
}

func (o *profileOperations) RegisterGet(api huma.API) {
	name := "Get Profile By ID"
	description := "Retrieve a specific profile by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"administrator"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the profile."`
	}) (*ProfileHolder, error) {
		profile, err := stores.NewProfileStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if profile == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &ProfileHolder{Body: profile}, nil
	})
}

func (o *profileOperations) RegisterGetAll(api huma.API) {
	name := "Get all Profiles"
	description := "Retrieve all profiles."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"administrator"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*ProfilesHolder, error) {
		profiles, err := stores.NewProfileStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &ProfilesHolder{Body: profiles}, nil
	})
}

func (o *profileOperations) RegisterGetMine(api huma.API) {
	name := "Get My Profile"
	description := "Retrieve the profile for the current user."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*ProfileHolder, error) {
		me, err := stores.NewProfileStore(relationalDB).Get(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		if me == nil {
			return nil, huma.Error404NotFound("The logged in user was not found.")
		}
		return &ProfileHolder{Body: me}, nil
	})
}
