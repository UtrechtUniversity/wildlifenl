package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewExperimentInput struct {
	Input
	Body *models.ExperimentRecord `json:"experiment"`
}

type UpdateExperimentInput struct {
	Input
	ID   string                   `query:"ID" format:"uuid" doc:"The ID of the experiment to be updated."`
	Body *models.ExperimentRecord `json:"experiment"`
}

type ExperimentHolder struct {
	Body *models.Experiment `json:"experiment"`
}

type ExperimentsHolder struct {
	Body []models.Experiment `json:"experiments"`
}

type experimentOperations Operations

func newExperimentOperations() *experimentOperations {
	return &experimentOperations{Endpoint: "experiment"}
}

/*
func (o *experimentOperations) RegisterGet(api huma.API) {
	name := "Get Experiment By ID"
	description := "Retrieve a specific experiment by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this experiment." format:"uuid"`
	}) (*ExperimentHolder, error) {
		experiment, err := stores.NewExperimentStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if experiment == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &ExperimentHolder{Body: experiment}, nil
	})
}
*/

func (o *experimentOperations) RegisterGetAll(api huma.API) {
	name := "Get All Experiments"
	description := "Retrieve all experiments."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*ExperimentsHolder, error) {
		experiments, err := stores.NewExperimentStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &ExperimentsHolder{Body: experiments}, nil
	})
}

func (o *experimentOperations) RegisterAdd(api huma.API) {
	name := "Add Experiment"
	description := "Add a new experiment."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewExperimentInput) (*ExperimentHolder, error) {
		experiment, err := stores.NewExperimentStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &ExperimentHolder{Body: experiment}, nil
	})
}

func (o *experimentOperations) RegisterGetMine(api huma.API) {
	name := "Get My Experiments"
	description := "Retrieve my experiments."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*ExperimentsHolder, error) {
		experiments, err := stores.NewExperimentStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ExperimentsHolder{Body: experiments}, nil
	})
}

func (o *experimentOperations) RegisterUpdate(api huma.API) {
	name := "Update Experiment"
	description := "Update an existing experiment that has not started yet."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *UpdateExperimentInput) (*ExperimentHolder, error) {
		experiment, err := stores.NewExperimentStore(relationalDB).Update(input.credential.UserID, input.ID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		if experiment == nil {
			return nil, generateNotFoundForThisUserError("experiment", input.ID)
		}
		return &ExperimentHolder{Body: experiment}, nil
	})
}
