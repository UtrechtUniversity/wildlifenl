package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type NoticeTypeStore Store

func NewNoticeTypeStore(db *sql.DB) *NoticeTypeStore {
	s := NoticeTypeStore{
		db: db,
		query: `
		SELECT t."id", t."nameNL", t."nameEN"
		FROM "noticeType" t
		`,
	}
	return &s
}

func (s *NoticeTypeStore) process(rows *sql.Rows, err error) ([]models.NoticeType, error) {
	if err != nil {
		return nil, err
	}
	noticeTypes := make([]models.NoticeType, 0)
	for rows.Next() {
		var noticeType models.NoticeType
		rows.Scan(&noticeType.ID, &noticeType.NameNL, &noticeType.NameEN)
		noticeTypes = append(noticeTypes, noticeType)
	}
	return noticeTypes, nil
}

func (s *NoticeTypeStore) GetAll() ([]models.NoticeType, error) {
	query := s.query + `
		ORDER BY t."id"
		`
	rows, err := s.db.Query(query)
	return s.process(rows, err)
}
