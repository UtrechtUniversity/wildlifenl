package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type AlarmHolder struct {
	Body *models.Alarm `json:"alarm"`
}

type AlarmsHolder struct {
	Body []models.Alarm `json:"alarms"`
}

type alarmOperations Operations

func newAlarmOperations() *alarmOperations {
	return &alarmOperations{Endpoint: "alarm"}
}

func (o *alarmOperations) RegisterGetAll(api huma.API) {
	name := "Get All Alarms"
	description := "Retrieve all alarms."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"administrator", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*AlarmsHolder, error) {
		alarms, err := stores.NewAlarmStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &AlarmsHolder{Body: alarms}, nil
	})
}

func (o *alarmOperations) RegisterGetMine(api huma.API) {
	name := "Get My Alarms"
	description := "Retrieve my alarms."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*AlarmsHolder, error) {
		alarms, err := stores.NewAlarmStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &AlarmsHolder{Body: alarms}, nil
	})
}
