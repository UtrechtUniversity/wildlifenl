package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewZoneInput struct {
	Input
	Body *models.ZoneRecord `json:"zone"`
}

type ZoneHolder struct {
	Body *models.Zone `json:"zone"`
}

type ZonesHolder struct {
	Body []models.Zone `json:"zones"`
}

type zoneOperations Operations

func newZoneOperations(database *sql.DB) *zoneOperations {
	o := zoneOperations{
		Database: database,
		Endpoint: "zone",
	}
	return &o
}

func (o *zoneOperations) RegisterGet(api huma.API) {
	name := "Get Zone By ID"
	description := "Retrieve a specific zone by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"administrator"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this zone." format:"uuid"`
	}) (*ZoneHolder, error) {
		zone, err := stores.NewZoneStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ZoneHolder{Body: zone}, nil
	})
}

func (o *zoneOperations) RegisterGetAll(api huma.API) {
	name := "Get All Zones"
	description := "Retrieve all zones."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"administrator", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*ZonesHolder, error) {
		zones, err := stores.NewZoneStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &ZonesHolder{Body: zones}, nil
	})
}

func (o *zoneOperations) RegisterAdd(api huma.API) {
	name := "Add Zone"
	description := "Add a new zone."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewZoneInput) (*ZoneHolder, error) {
		species, err := stores.NewZoneStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &ZoneHolder{Body: species}, nil
	})
}

func (o *zoneOperations) RegisterGetMine(api huma.API) {
	name := "Get My Zones"
	description := "Retrieve my zones."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*ZonesHolder, error) {
		zones, err := stores.NewZoneStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ZonesHolder{Body: zones}, nil
	})
}
