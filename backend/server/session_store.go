package server

import (
	"database/sql"

	"github.com/guricerin/grumbler/backend/model"
)

type sessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) sessionStore {
	return sessionStore{db: db}
}

func (s *sessionStore) Create(token string, user model.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("insert into sessions (token, user_pk) values (?, ?)", token, user.Pk)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *sessionStore) RetrieveByToken(token string) (model.Session, error) {
	sess := model.Session{}
	tx, err := s.db.Begin()
	if err != nil {
		return sess, err
	}
	err = tx.QueryRow("select pk, token, user_pk from sessions where token = ?", token).Scan(&sess.Pk, &sess.Token, &sess.UserPk)
	if err != nil {
		tx.Rollback()
		return sess, err
	}
	err = tx.Commit()
	return sess, err
}

func (s *sessionStore) Update(oldToken string, newToken string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("update sessions set token = ? where token = ?", newToken, oldToken)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *sessionStore) DeleteByToken(token string) error {
	_, err := s.db.Exec("delete from sessions where token = ?", token)
	return err
}

func (s *sessionStore) DeleteByUserPk(userPk uint64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from sessions where user_pk = ?", userPk)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
