package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type ContactStartInput struct {
	Input
	Body *struct {
		ContactHardwareAddress string        `json:"contactHardwareAddress" pattern:"^([0-9A-F]{2}:){5}[0-9A-F]{2}$" doc:"The EUI-48 hardware address of the bluetooth contact tracing device."`
		Location               *models.Point `json:"location,omitempty" doc:"The location of the user at the moment this contact started."`
	}
}

type ContactEndInput struct {
	Input
	Body *struct {
		ContactID string        `json:"contactID" format:"uuid" doc:"The ID of the contact tracing event to end."`
		Location  *models.Point `json:"location,omitempty" doc:"The location of the user at the moment this contact finished."`
	}
}

type ContactHolder struct {
	Body *models.Contact `json:"contact"`
}

type ContactsHolder struct {
	Body []models.Contact `json:"contacts"`
}

type contactOperations Operations

func newContactOperations() *contactOperations {
	return &contactOperations{Endpoint: "contact"}
}

func (o *contactOperations) RegisterStart(api huma.API) {
	name := "Start Contact"
	description := "Start a new contact tracing event."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *ContactStartInput) (*ContactHolder, error) {
		store := stores.NewContactStore(relationalDB)
		ok, err := store.NotExists(input.credential.UserID, input.Body.ContactHardwareAddress)
		if err != nil {
			return nil, handleError(err)
		}
		if !ok {
			return nil, huma.Error409Conflict("An existing active contact tracing record exists for the given contactHardwareAddress.")
		}
		contact, err := store.Add(input.credential.UserID, input.Body.ContactHardwareAddress, input.Body.Location)
		if err != nil {
			return nil, handleError(err)
		}
		if contact == nil {
			return nil, huma.Error404NotFound("A borne-sensor deployment was not found for the provided contactHardwareAddress.")
		}

		// TODO: Add conveyances here.

		return &ContactHolder{Body: contact}, nil
	})
}

func (o *contactOperations) RegisterEnd(api huma.API) {
	name := "End Contact"
	description := "End an existing contact tracing event."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *ContactEndInput) (*ContactHolder, error) {
		store := stores.NewContactStore(relationalDB)
		ok, err := store.Exists(input.credential.UserID, input.Body.ContactID)
		if err != nil {
			return nil, handleError(err)
		}
		if !ok {
			return nil, huma.Error409Conflict("An existing active contact tracing record was not found.")
		}
		contact, err := store.End(input.credential.UserID, input.Body.ContactID, input.Body.Location)
		if err != nil {
			return nil, handleError(err)
		}
		return &ContactHolder{Body: contact}, nil
	})
}

func (o *contactOperations) RegisterGetMine(api huma.API) {
	name := "Get My Contacts"
	description := "Retrieve my contact tracing events."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*ContactsHolder, error) {
		contacts, err := stores.NewContactStore(relationalDB).GetForUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ContactsHolder{Body: contacts}, nil
	})
}

func (o *contactOperations) RegisterGetByExperiment(api huma.API) {
	name := "Get All Contacts By Experiment"
	description := "Retrieve all contact tracing events that have a spatiotemporal overlap with an experiment. A contact overlaps temporally if its start or end falls within the timespan of the experiment, and spatially if the location of the user at its start or end is within the living lab of the experiment. If the experiment is not bound to a living lab, every contact has a spatial overlap with it."
	path := "/" + o.Endpoint + "s/experiment/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the experiment to retrieve contacts for."`
	}) (*ContactsHolder, error) {
		contacts, err := stores.NewContactStore(relationalDB).GetByExperiment(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &ContactsHolder{Body: contacts}, nil
	})
}
