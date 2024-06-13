package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewNoticeInput struct {
	Input
	Body *models.NoticeRecord `json:"notice"`
}

type NoticeHolder struct {
	Body *models.Notice `json:"notice"`
}

type NoticesHolder struct {
	Body []models.Notice `json:"notices"`
}

type noticeOperations Operations

func newNoticeOperations(database *sql.DB) *noticeOperations {
	o := noticeOperations{
		Database: database,
		Endpoint: "notice",
	}
	return &o
}

func (o *noticeOperations) RegisterGet(api huma.API) {
	name := "Get Notice By ID"
	description := "Retrieve a specific notice by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"wildlife-manager", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the notice."`
	}) (*NoticeHolder, error) {
		notice, err := stores.NewNoticeStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if notice == nil {
			return nil, huma.Error404NotFound("No notice with ID " + input.ID + " was found")
		}
		return &NoticeHolder{Body: notice}, nil
	})
}

func (o *noticeOperations) RegisterGetAll(api huma.API) {
	name := "Get All Notices"
	description := "Retrieve all notices."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"wildlife-manager", "researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*NoticesHolder, error) {
		notices, err := stores.NewNoticeStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &NoticesHolder{Body: notices}, nil
	})
}

func (o *noticeOperations) RegisterGetMy(api huma.API) {
	name := "Get My Notices"
	description := "Retrieve all notices made by the current user."
	path := "/" + o.Endpoint + "/me/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*NoticesHolder, error) {
		notices, err := stores.NewNoticeStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &NoticesHolder{Body: notices}, nil
	})
}

func (o *noticeOperations) RegisterAdd(api huma.API) {
	name := "Add New Notice"
	description := "Submit a new notice."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewNoticeInput) (*NoticeHolder, error) {
		notice, err := stores.NewNoticeStore(relationalDB).Add(input.credential.UserID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &NoticeHolder{Body: notice}, nil
	})
}
