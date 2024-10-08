package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type AddRoleToUserInput struct {
	Body *struct {
		UserID string `json:"userID" format:"uuid" doc:"The ID of the user"`
		RoleID int    `json:"roleID" minimum:"1" doc:"The ID of the role"`
	}
}

type RolesHolder struct {
	Body []models.Role `json:"roles"`
}

type roleOperations Operations

func newRoleOperations() *roleOperations {
	return &roleOperations{Endpoint: "role"}
}

func (o *roleOperations) RegisterGetAll(api huma.API) {
	name := "Get All Roles"
	description := "Retrieve all roles."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*RolesHolder, error) {
		roles, err := stores.NewRoleStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &RolesHolder{Body: roles}, nil
	})
}

func (o *roleOperations) RegisterAddRoleToUser(api huma.API) {
	name := "Add a Role to a User"
	description := "Add a specific role to a specific user."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *AddRoleToUserInput) (*struct{}, error) {
		err := stores.NewRoleStore(relationalDB).AddRoleToUser(input.Body.UserID, input.Body.RoleID)
		if err != nil {
			return nil, handleError(err)
		}
		return nil, nil
	})
}
