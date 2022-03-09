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
    from grumbles as g
    left join users as u
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

func (s *grumbleStore) CreateBookmark(grumblePk string, byUserId string) (model.Bookmark, error) {
	query := `insert into bookmarks
    (grumble_pk, by_user_id)
    values (?, ?)`
	_, err := s.db.Exec(query, grumblePk, byUserId)

	res := model.Bookmark{
		GrumblePk: grumblePk,
		ByUserId:  byUserId,
	}
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *grumbleStore) retrieveBookmarksByUserId(userId string) ([]model.Bookmark, error) {
	res := make([]model.Bookmark, 0)
	query := `select pk, grumble_pk, by_user_id from bookmarks
    where by_user_id = ?`
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		b := model.Bookmark{}
		err = rows.Scan(&b.Pk, &b.GrumblePk, &b.ByUserId)
		if err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, nil
}

func (s *grumbleStore) RetrieveBookmarkedGrumblesByUserId(userId string) ([]model.GrumbleRes, error) {
	bookmarks, err := s.retrieveBookmarksByUserId(userId)
	if err != nil {
		return nil, err
	}

	res := make([]model.GrumbleRes, 0)
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    where g.pk = ?`
	for _, b := range bookmarks {
		rows, err := s.db.Query(query, b.GrumblePk)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// ループは一周しかしない（はず）
		for rows.Next() {
			g := model.GrumbleRes{}
			err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
			res = append(res, g)
		}
	}

	return res, nil
}
