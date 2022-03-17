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

func (s *userStore) DeleteByPk(pk uint64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	// 対象ユーザに対応する他のテーブルの行も削除する
	query := `delete u, s, g, b, reg, f, r from users as u
    left join sessions as s
        on u.pk = s.user_pk
    left join grumbles as g
        on u.id = g.user_id
    left join bookmarks as b
        on u.id = b.by_user_id
    left join regrumbles as reg
        on u.id = reg.by_user_id or u.id = reg.dst_user_id
    left join follows as f
        on u.id = f.src_user_id or u.id = f.dst_user_id
    left join replies as r
        on g.pk = r.dst_grumble_pk or g.pk = r.src_grumble_pk
    where u.pk = ?`
	_, err = tx.Exec(query, pk)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

const (
	UserIdSearch = iota
	UserNameSearch
)

func (s *userStore) Search(keyword string, kind int) ([]model.User, error) {
	var query string
	switch kind {
	case UserIdSearch:
		query = `select pk, id, name, password, profile from users
        where id like concat('%', ?, '%')`
	case UserNameSearch:
		query = `select pk, id, name, password, profile from users
        where name like concat('%', ?, '%')`
	}
	return s.searchInner(query, keyword)
}

func (s *userStore) searchInner(query string, keyword string) ([]model.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	users := make([]model.User, 0)
	rows, err := tx.Query(query, keyword)
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

	return users, tx.Commit()
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
