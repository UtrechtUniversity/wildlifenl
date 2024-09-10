package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewBorneSensorDeploymentInput struct {
	Body *models.BorneSensorDeploymentRecord `json:"borneSensorDeployment"`
}

type BorneSensorDeploymentHolder struct {
	Body *models.BorneSensorDeployment `json:"borneSensorDeployment"`
}

type BorneSensorDeploymentsHolder struct {
	Body []models.BorneSensorDeployment `json:"borneSensorDeployments"`
}

type borneSensorDeploymentOperations Operations

func newBorneSensorDeploymentOperations() *borneSensorDeploymentOperations {
	o := borneSensorDeploymentOperations{
		Endpoint: "borne-sensor-deployment",
	}
	return &o
}

func (o *borneSensorDeploymentOperations) RegisterGetAll(api huma.API) {
	name := "Get All BorneSensorDeployments"
	description := "Retrieve all active borne sensor deployments, i.e. having no end timestamp."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"herd-manager"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*BorneSensorDeploymentsHolder, error) {
		borneSensorDeployments, err := stores.NewBorneSensorDeploymentStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorDeploymentsHolder{Body: borneSensorDeployments}, nil
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
	}, func(ctx context.Context, input *NewBorneSensorDeploymentInput) (*BorneSensorDeploymentHolder, error) {
		borneSensorDeployment, err := stores.NewBorneSensorDeploymentStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &BorneSensorDeploymentHolder{Body: borneSensorDeployment}, nil
	})
}
