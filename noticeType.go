package wildlifenl

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
)

type NoticeTypesHolder struct {
	Body []models.NoticeType `json:"noticeTypes"`
}

type noticeTypeOperations struct{}

func (s *noticeOperations) RegisterGetNoticeTypes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-noticeTypes",
		Tags:        []string{"noticeType"},
		Method:      http.MethodGet,
		Path:        "/noticetype/",
		Summary:     "Get all NoticeTypes",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve all noticeTypes available.",
	}, func(ctx context.Context, input *struct{}) (*NoticeTypesHolder, error) {
		noticeTypes, err := stores.NewNoticeTypeStore(database).GetAll()
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &NoticeTypesHolder{Body: noticeTypes}, nil
	})
}
