package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type UserStore Store

func NewUserStore(db *sql.DB) *UserStore {
	s := UserStore{
		relationalDB: db,
		query: `
		SELECT u."id", u."name", r."id", r."name"
		FROM "user" u
		LEFT JOIN user_role x ON x."userID" = u."id"
		LEFT JOIN role r ON r."id" = x."roleID"
		`,
	}
	return &s
}

func (s *UserStore) process(rows *sql.Rows, err error) ([]models.User, error) {
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	var user models.User
	for rows.Next() {
		var userID string
		var userName string
		var role models.Role
		if err := rows.Scan(&userID, &userName, &role.ID, &role.Name); err != nil {
			return nil, err
		}
		if user.ID != "" && user.ID != userID {
			users = append(users, user)
			user = models.User{}
		}
		user.ID = userID
		user.Name = userName
		if role.ID > 0 {
			user.Roles = append(user.Roles, role)
		}
	}
	if user.ID != "" {
		users = append(users, user)
	}
	return users, nil
}

func (s *UserStore) Get(userID string) (*models.User, error) {
	query := s.query + `
		WHERE u."id" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *UserStore) GetAll() ([]models.User, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}
