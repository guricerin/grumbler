package server

import (
	"database/sql"
	"time"

	"github.com/guricerin/grumbler/backend/model"
)

type grumbleStore struct {
	db *sql.DB
}

func NewGrumbleStore(db *sql.DB) grumbleStore {
	return grumbleStore{db: db}
}

func (s *grumbleStore) Create(content string, user model.User) error {
	t := time.Now()
	id, err := createUlid(t)
	if err != nil {
		return err
	}

	g := model.Grumble{
		Pk:        id,
		Content:   content,
		UserId:    user.Id,
		CreatedAt: t,
	}
	_, err = s.db.Exec("insert into grumbles (pk, content, user_id, created_at) values (?, ?, ?, ?)", g.Pk, g.Content, g.UserId, g.CreatedAt)
	return err
}

func (s *grumbleStore) RetrieveByUserId(userId string) ([]model.Grumble, error) {
	res := make([]model.Grumble, 0)
	rows, err := s.db.Query("select pk, content, user_id, created_at from grumbles where user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		g := model.Grumble{}
		err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, g)
	}
	return res, nil
}

func (s *grumbleStore) DeleteByPk(pk string) error {
	_, err := s.db.Exec("delete grumbles where pk = ?", pk)
	return err
}
