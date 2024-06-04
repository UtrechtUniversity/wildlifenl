package wildlifenl

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type MyNoticesInput struct {
	Input
}

type NewNoticeInput struct {
	Input
	Body *models.Notice `json:"notice"`
}

type NoticeHolder struct {
	Body *models.Notice `json:"notice"`
}

type NoticesHolder struct {
	Body []models.Notice `json:"notices"`
}

type noticeOperations struct{}

func (s *noticeOperations) RegisterGet(api huma.API) {
	scopes := []string{"wildlife-manager", "researcher"}
	huma.Register(api, huma.Operation{
		OperationID: "get-notice-by-id",
		Tags:        []string{"notice"},
		Method:      http.MethodGet,
		Path:        "/notice/{id}",
		Summary:     "Get Notice By ID",
		Security:    []map[string][]string{{"auth": scopes}},
		Description: "Retrieve a specific notice by its ID. <br/><br/>**Scopes**<br/>" + scopesAsMarkdown(scopes),
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the notice." format:"uuid"`
	}) (*NoticeHolder, error) {
		notice, err := stores.NewNoticeStore(database).Get(input.ID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		if notice == nil {
			return nil, huma.Error404NotFound("No notice with ID " + input.ID + " was found")
		}
		return &NoticeHolder{Body: notice}, nil
	})
}

func (s *noticeOperations) RegisterGetAll(api huma.API) {
	scopes := []string{"wildlife-manager", "researcher"}
	huma.Register(api, huma.Operation{
		OperationID: "get-all-notices",
		Tags:        []string{"notice"},
		Method:      http.MethodGet,
		Path:        "/notice/",
		Summary:     "Get all Notices",
		Security:    []map[string][]string{{"auth": scopes}},
		Description: "Retrieve all notices. <br/><br/>**Scopes**<br/>" + scopesAsMarkdown(scopes),
	}, func(ctx context.Context, input *struct{}) (*NoticesHolder, error) {
		notices, err := stores.NewNoticeStore(database).GetAll()
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &NoticesHolder{Body: notices}, nil
	})
}

func (s *noticeOperations) RegisterGetMy(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-my-notices",
		Tags:        []string{"notice"},
		Method:      http.MethodGet,
		Path:        "/notice/me/",
		Summary:     "Get my Notices",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve all notices made by the current user.",
	}, func(ctx context.Context, input *MyNoticesInput) (*NoticesHolder, error) {
		notices, err := stores.NewNoticeStore(database).GetByUser(input.credential.UserID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &NoticesHolder{Body: notices}, nil
	})
}

func (s *noticeOperations) RegisterAdd(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "add-notice",
		Tags:        []string{"notice"},
		Method:      http.MethodPost,
		Path:        "/notice/",
		Summary:     "Add a new Notice",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Submit a new notice.",
	}, func(ctx context.Context, input *NewNoticeInput) (*NoticeHolder, error) {
		noticeTypes, err := stores.NewNoticeTypeStore(database).GetAll()
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		noticeTypeFound := false
		for _, noticeType := range noticeTypes {
			if noticeType.ID == input.Body.Type.ID {
				noticeTypeFound = true
				break
			}
		}
		if !noticeTypeFound {
			return nil, huma.Error400BadRequest("NoticeType ID " + strconv.Itoa(input.Body.Type.ID) + " is not valid.")
		}
		notice, err := stores.NewNoticeStore(database).Add(input.credential.UserID, input.Body)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &NoticeHolder{Body: notice}, nil
	})
}
