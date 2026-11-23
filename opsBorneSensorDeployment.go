package wildlifenl

import (
	"context"
	"net/http"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type BorneSensorDeploymentHolder struct {
	Body *models.BorneSensorDeployment `json:"borneSensorDeployment"`
}

type BorneSensorDeploymentsHolder struct {
	Body []models.BorneSensorDeployment `json:"borneSensorDeployments"`
}

type BorneSensorDeploymentAddInput struct {
	Input
	Body *models.BorneSensorDeploymentRecord `json:"borneSensorDeployment"`
}

type BorneSensorDeploymentEndInput struct {
	Input
	Body *struct {
		ID        string    `json:"ID" format:"uuid" doc:"The ID of the borne-sensor deployment to be ended."`
		TimeStamp time.Time `json:"timestamp" doc:"The moment the borne-sensor deployment ended."`
	}
}

type borneSensorDeploymentOperations Operations

func newBorneSensorDeploymentOperations() *borneSensorDeploymentOperations {
	o := borneSensorDeploymentOperations{
		Endpoint: "borne-sensor-deployment",
	}
	return &o
}

func (o *borneSensorDeploymentOperations) RegisterGetMine(api huma.API) {
	name := "Get My BorneSensorDeployments"
	description := "Get all borne-sensor deployments for the current user."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{"herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*BorneSensorDeploymentsHolder, error) {
		borneSensorDeployment, err := stores.NewBorneSensorDeploymentStore(relationalDB, timeseriesDB).GetMine(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorDeploymentsHolder{Body: borneSensorDeployment}, nil
	})
}

func (o *borneSensorDeploymentOperations) RegisterGetMyActive(api huma.API) {
	name := "Get My Active BorneSensorDeployments"
	description := "Get all active borne-sensor deployments for the current user."
	path := "/" + o.Endpoint + "/me/active/"
	scopes := []string{"herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*BorneSensorDeploymentsHolder, error) {
		borneSensorDeployment, err := stores.NewBorneSensorDeploymentStore(relationalDB, timeseriesDB).GetMyActive(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorDeploymentsHolder{Body: borneSensorDeployment}, nil
	})
}

func (o *borneSensorDeploymentOperations) RegisterAdd(api huma.API) {
	name := "Add BorneSensorDeployment"
	description := "Submit a new deployment for a borne sensor. If an existing deployment is found for the same sensorID, that deployment's end timestamp will be set to the start timestamp of the new deployment."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"herd-manager"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *BorneSensorDeploymentAddInput) (*BorneSensorDeploymentHolder, error) {
		userID := input.credential.UserID
		borneSensorDeployment, err := stores.NewBorneSensorDeploymentStore(relationalDB, timeseriesDB).Add(userID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorDeploymentHolder{Body: borneSensorDeployment}, nil
	})
}

func (o *borneSensorDeploymentOperations) RegisterEnd(api huma.API) {
	name := "End BorneSensorDeployment"
	description := "End a deployment for a borne sensor."
	path := "/" + o.Endpoint + "/end"
	scopes := []string{"herd-manager"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *BorneSensorDeploymentEndInput) (*BorneSensorDeploymentHolder, error) {
		borneSensorDeployment, err := stores.NewBorneSensorDeploymentStore(relationalDB, timeseriesDB).End(input.credential.UserID, input.Body.ID, input.Body.TimeStamp)
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorDeploymentHolder{Body: borneSensorDeployment}, nil
	})
}
