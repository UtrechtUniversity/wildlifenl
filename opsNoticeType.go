package wildlifenl

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NoticeTypesHolder struct {
	Body []models.NoticeType `json:"noticeTypes"`
}

type noticeTypeOperations Operations

func newNoticeTypeOperations(database *sql.DB) *noticeTypeOperations {
	o := noticeTypeOperations{
		Database: database,
		Endpoint: "notice-type",
	}
	return &o
}

func (o *noticeTypeOperations) RegisterGetAll(api huma.API) {
	name := "Get All NoticeTypes"
	description := "Retrieve all notice types."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*NoticeTypesHolder, error) {
		noticeTypes, err := stores.NewNoticeTypeStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &NoticeTypesHolder{Body: noticeTypes}, nil
	})
}
