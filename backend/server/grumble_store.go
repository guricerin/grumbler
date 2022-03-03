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

func (s *grumbleStore) RetrieveByUserId(userId string) ([]model.GrumbleRes, error) {
	res := make([]model.GrumbleRes, 0)
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g left join users as u
    on g.user_id = u.id
    where u.id = ?`
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		g := model.GrumbleRes{}
		err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
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
