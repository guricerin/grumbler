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

func (s *userStore) Create(user *model.User) error {
	res, err := s.db.Exec("insert into users (id, name, password, profile) values (?, ?, ?, ?)", user.Id, user.Name, user.Password, user.Profile)
	if err != nil {
		return err
	}
	pk, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.Pk = uint64(pk)
	return err
}

func (s *userStore) RetrieveById(userId string) (model.User, error) {
	user := model.User{}
	err := s.db.QueryRow("select pk, id, name, password, profile from users where id = ?", userId).Scan(&user.Pk, &user.Id, &user.Name, &user.Password, &user.Profile)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userStore) RetrieveByPk(pk uint64) (model.User, error) {
	user := model.User{}
	err := s.db.QueryRow("select pk, id, name, password, profile from users where pk = ?", pk).Scan(&user.Pk, &user.Id, &user.Name, &user.Password, &user.Profile)
	if err != nil {
		return user, err
	}
	return user, nil
}

// 対象ユーザに対応する他のテーブルの行も削除する
func (s *userStore) DeleteByPk(pk uint64) error {
	_, err := s.db.Exec("delete u, s from users as u left join sessions as s on u.pk = s.user_pk where u.pk = ?", pk)
	return err
}

func (s *userStore) SearchById(id string) ([]model.User, error) {
	users := make([]model.User, 0)
	pattern := "%" + id + "%"
	// rows, err := s.db.Query("select pk, id, name, password from users where id like '%' || ? || '%'", id)
	rows, err := s.db.Query("select pk, id, name, password, profile from users where id like ?", pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.Pk, &user.Id, &user.Name, &user.Password, &user.Profile)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
