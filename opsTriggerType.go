package wildlifenl

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type TriggerTypesHolder struct {
	Body []string `json:"triggerTypes"`
}

type triggerTypeOperations Operations

func newTriggerTypeOperations() *triggerTypeOperations {
	return &triggerTypeOperations{Endpoint: "triggerType"}
}

func (o *triggerTypeOperations) RegisterGetAll(api huma.API) {
	name := "Get All TriggerTypes"
	description := "Retrieve all trigger types."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*TriggerTypesHolder, error) {
		// This is not super elegant and hard coded, but a client needed it.
		return &TriggerTypesHolder{Body: []string{"encounter", "answer", "alarm"}}, nil
	})
}
