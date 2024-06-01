package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type NoticeStore Store

func NewNoticeStore(db *sql.DB) *NoticeStore {
	s := NoticeStore{
		db: db,
		query: `
		SELECT n."id", n."description", n."latitude", n."longitude", t."id", t."nameNL", t."nameEN", u."id", u."name"
		FROM notice n
		INNER JOIN "noticeType"	t ON t."id" = n."typeID"
		INNER JOIN "user" u ON u."id" = n."userID"
		`,
	}
	return &s
}

func (s *NoticeStore) process(rows *sql.Rows, err error) ([]models.Notice, error) {
	if err != nil {
		return nil, err
	}
	notices := make([]models.Notice, 0)
	for rows.Next() {
		var notice models.Notice
		var noticeType models.NoticeType
		var user models.User
		rows.Scan(&notice.ID, &notice.Description, &notice.Latitude, &notice.Longitude, &noticeType.ID, &noticeType.NameNL, &noticeType.NameEN, &user.ID, &user.Name)
		notice.Type = noticeType
		notice.Reporter = user
		notices = append(notices, notice)
	}
	return notices, nil
}

func (s *NoticeStore) Get(noticeID string) (*models.Notice, error) {
	query := s.query + `
		WHERE n."id" = $1
		`
	rows, err := s.db.Query(query, noticeID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *NoticeStore) GetAll() ([]models.Notice, error) {
	rows, err := s.db.Query(s.query)
	return s.process(rows, err)
}

func (s *NoticeStore) GetByUser(userID string) ([]models.Notice, error) {
	query := s.query + `
		WHERE u."id" = $1
		`
	rows, err := s.db.Query(query, userID)
	return s.process(rows, err)
}

func (s *NoticeStore) Add(userID string, notice *models.Notice) (*models.Notice, error) {
	query := `
		INSERT INTO notice("description", "latitude", "longitude", "typeID", "userID") VALUES($1, $2, $3, $4, $5)
		RETURNING "id"
	`
	var id string
	row := s.db.QueryRow(query, notice.Description, notice.Latitude, notice.Longitude, notice.Type.ID, userID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
