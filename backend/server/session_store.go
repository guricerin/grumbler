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
	_, err := s.db.Exec("insert into sessions (token, user_pk) values (?, ?)", token, user.Pk)
	return err
}

func (s *sessionStore) RetrieveByToken(token string) (model.Session, error) {
	sess := model.Session{}
	err := s.db.QueryRow("select pk, token, user_pk from sessions where token = ?", token).Scan(&sess.Pk, &sess.Token, &sess.UserPk)
	if err != nil {
		return sess, err
	}
	return sess, nil
}

func (s *sessionStore) Update(oldToken string, newToken string) error {
	_, err := s.db.Exec("update sessions set token = ? where token = ?", newToken, oldToken)
	return err
}

func (s *sessionStore) DeleteByToken(token string) error {
	_, err := s.db.Exec("delete from sessions where token = ?", token)
	return err
}

func (s *sessionStore) DeleteByUserPk(userPk uint64) error {
	_, err := s.db.Exec("delete from sessions where user_pk = ?", userPk)
	return err
}
