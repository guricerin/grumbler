package server

import (
	"database/sql"
	"time"

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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("insert into users (id, created_at, name, password, profile) values (?, ?, ?, ?, ?)", user.Id, time.Now(), user.Name, user.Password, user.Profile)
	if err != nil {
		tx.Rollback()
		return err
	}
	pk, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	user.Pk = uint64(pk)
	err = tx.Commit()
	return err
}

func (s *userStore) RetrieveById(userId string) (model.User, error) {
	user := model.User{}
	tx, err := s.db.Begin()
	if err != nil {
		return user, err
	}
	err = tx.QueryRow("select pk, id, name, password, profile from users where id = ?", userId).Scan(&user.Pk, &user.Id, &user.Name, &user.Password, &user.Profile)
	if err != nil {
		tx.Rollback()
		return user, err
	}
	err = tx.Commit()
	return user, err
}

func (s *userStore) RetrieveByPk(pk uint64) (model.User, error) {
	user := model.User{}
	tx, err := s.db.Begin()
	if err != nil {
		return user, err
	}
	err = tx.QueryRow("select pk, id, name, password, profile from users where pk = ?", pk).Scan(&user.Pk, &user.Id, &user.Name, &user.Password, &user.Profile)
	if err != nil {
		tx.Rollback()
		return user, err
	}
	err = tx.Commit()
	return user, err
}

// todo: 対象ユーザに対応する他のテーブルの行も削除する
func (s *userStore) DeleteByPk(pk uint64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete u, s from users as u left join sessions as s on u.pk = s.user_pk where u.pk = ?", pk)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *userStore) SearchById(id string) ([]model.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	users := make([]model.User, 0)
	pattern := "%" + id + "%"
	// rows, err := s.db.Query("select pk, id, name, password from users where id like '%' || ? || '%'", id)
	rows, err := tx.Query("select pk, id, name, password, profile from users where id like ?", pattern)
	if err != nil {
		tx.Rollback()
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

	err = tx.Commit()
	return users, err
}

func (s *userStore) SearchByName(name string) ([]model.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	users := make([]model.User, 0)
	pattern := "%" + name + "%"
	query := `select pk, id, name, password, profile from users
    where name like ?`
	rows, err := tx.Query(query, pattern)
	if err != nil {
		tx.Rollback()
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

	err = tx.Commit()
	return users, err
}

func (s *userStore) Update(user *model.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `update users
    set name = ?, profile = ?
    where id = ?`
	_, err = tx.Exec(query, user.Name, user.Profile, user.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
