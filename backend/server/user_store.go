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

func (s *userStore) RetrieveAllById(userId string) ([]model.User, error) {
	query := `select pk, id, name, password, profile from users
    where id = ?`
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]model.User, 0)
	for rows.Next() {
		u := model.User{}
		err := rows.Scan(&u.Pk, &u.Id, &u.Name, &u.Password, &u.Profile)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
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

func (s *userStore) SearchByName(name string) ([]model.User, error) {
	users := make([]model.User, 0)
	pattern := "%" + name + "%"
	query := `select pk, id, name, password, profile from users
    where name like ?`
	rows, err := s.db.Query(query, pattern)
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

func (s *userStore) Update(user *model.User) error {
	query := `update users
    set name = ?, profile = ?
    where id = ?`
	_, err := s.db.Exec(query, user.Name, user.Profile, user.Id)
	return err
}
