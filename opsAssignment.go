package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type AssignmentsHolder struct {
	Body []models.Assignment `json:"assignment"`
}

type assignmentOperations Operations

func newAssignmentOperations() *assignmentOperations {
	return &assignmentOperations{Endpoint: "assigment"}
}

func (o *assignmentOperations) RegisterGetMine(api huma.API) {
	name := "Get My Assignments"
	description := "Retrieve the assignments for the current user."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*AssignmentsHolder, error) {
		assignments, err := stores.NewAssignmentStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &AssignmentsHolder{Body: assignments}, nil
	})
}
