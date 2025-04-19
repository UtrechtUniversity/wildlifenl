package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type ZoneHolder struct {
	Body *models.Zone `json:"zone"`
}

type ZonesHolder struct {
	Body []models.Zone `json:"zones"`
}

type ZoneAddInput struct {
	Input
	Body *models.ZoneRecord `json:"zone"`
}

type SpeciesForZoneUpdateInput struct {
	Input
	Body *struct {
		ZoneID    string `json:"zoneID" format:"uuid" doc:"The ID of the zone"`
		SpeciesID string `json:"speciesID" format:"uuid" doc:"The ID of the species"`
	}
}

type ZoneDeactivateInput struct {
	Input
	ID string `path:"id" doc:"The ID of this zone." format:"uuid"`
}

type zoneOperations Operations

func newZoneOperations() *zoneOperations {
	return &zoneOperations{Endpoint: "zone"}
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
	path := "/" + o.Endpoint + "s/"
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
	}, func(ctx context.Context, input *ZoneAddInput) (*ZoneHolder, error) {
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
	path := "/" + o.Endpoint + "s/me/"
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

func (o *zoneOperations) RegisterAddSpeciesToZone(api huma.API) {
	name := "Add a Species to a Zone"
	description := "Add a species for which this zone should create alarms."
	path := "/" + o.Endpoint + "/species/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *SpeciesForZoneUpdateInput) (*ZoneHolder, error) {
		zone, err := stores.NewZoneStore(relationalDB).AddSpeciesToZone(input.credential.UserID, input.Body.ZoneID, input.Body.SpeciesID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ZoneHolder{Body: zone}, nil
	})
}

func (o *zoneOperations) RegisterRemoveSpeciesFromZone(api huma.API) {
	name := "Remove a Species from a Zone"
	description := "Remove a species for which this zone should create alarms."
	path := "/" + o.Endpoint + "/species/"
	scopes := []string{}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *SpeciesForZoneUpdateInput) (*ZoneHolder, error) {
		zone, err := stores.NewZoneStore(relationalDB).RemoveSpeciesFromZone(input.credential.UserID, input.Body.ZoneID, input.Body.SpeciesID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ZoneHolder{Body: zone}, nil
	})
}

func (o *zoneOperations) RegisterDeactivate(api huma.API) {
	name := "Deactivate Zone"
	description := "Deactivate a zone."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodDelete
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *ZoneDeactivateInput) (*ZoneHolder, error) {
		store := stores.NewZoneStore(relationalDB)
		var zone models.Zone
		zones, err := store.GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		for _, z := range zones {
			if z.ID == input.ID {
				zone = z
				break
			}
		}
		if zone.ID == "" {
			return nil, generateNotFoundForThisUserError(o.Endpoint, input.ID)
		}
		result, err := stores.NewZoneStore(relationalDB).Deactivate(zone.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ZoneHolder{Body: result}, nil
	})
}
