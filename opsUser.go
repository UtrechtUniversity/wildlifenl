package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type UserHolder struct {
	Body *models.User `json:"user"`
}

type UsersHolder struct {
	Body []models.User `json:"users"`
}

type userOperations Operations

func newUserOperations(database *sql.DB) *userOperations {
	o := userOperations{
		Database: database,
		Endpoint: "user",
	}
	return &o
}

func (o *userOperations) RegisterGetUserByID(api huma.API) {
	name := "Get User By ID"
	description := "Retrieve a specific user by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the user."`
	}) (*UserHolder, error) {
		user, err := stores.NewUserStore(database).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if user == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &UserHolder{Body: user}, nil
	})
}

func (o *userOperations) RegisterGetAllUsers(api huma.API) {
	name := "Get all Users"
	description := "Retrieve all users."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*UsersHolder, error) {
		users, err := stores.NewUserStore(database).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &UsersHolder{Body: users}, nil
	})
}
