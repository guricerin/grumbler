package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guricerin/grumbler/backend/model"
)

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) userStore {
	return userStore{db: db}
}

func (s *userStore) Create(user model.User) error {
	_, err := s.db.Exec("insert into users (id, name, password) values (?, ?, ?)", user.Id, user.Name, user.Password)
	return err
}

func (s *userStore) RetrieveById(userId string) (model.User, error) {
	user := model.User{}
	err := s.db.QueryRow("select pk, id, name, password from users where id = ?", userId).Scan(&user.Pk, &user.Id, &user.Name, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userStore) RetrieveByPk(pk uint) (model.User, error) {
	user := model.User{}
	err := s.db.QueryRow("select pk, id, name, password from users where pk = ?", pk).Scan(&user.Pk, &user.Id, &user.Name, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
