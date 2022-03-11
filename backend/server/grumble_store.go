package server

import (
	"database/sql"
	"errors"
	"time"

	"github.com/guricerin/grumbler/backend/model"
)

type grumbleStore struct {
	db *sql.DB
}

func NewGrumbleStore(db *sql.DB) grumbleStore {
	return grumbleStore{db: db}
}

func (s *grumbleStore) RetrieveByPk(grumblePk string, signinUserId string) (model.GrumbleRes, error) {
	res := model.GrumbleRes{}
	tx, err := s.db.Begin()
	if err != nil {
		return res, err
	}
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    where g.pk = ?`
	row := tx.QueryRow(query, grumblePk)
	err = row.Scan(&res.Pk, &res.Content, &res.UserId, &res.CreatedAt, &res.UserName)
	if err != nil {
		return res, err
	}
	err = s.retrieveBookmarkedCountAndBySigninUser(tx, &res, signinUserId)
	if err != nil {
		tx.Rollback()
		return res, err
	}
	return res, tx.Commit()
}

func (s *grumbleStore) Create(content string, user model.User) error {
	t := time.Now()
	id, err := createUlid(t)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	g := model.Grumble{
		Pk:        id,
		Content:   content,
		UserId:    user.Id,
		CreatedAt: t,
	}
	_, err = tx.Exec("insert into grumbles (pk, content, user_id, created_at) values (?, ?, ?, ?)", g.Pk, g.Content, g.UserId, g.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *grumbleStore) retrieveBookmarkedCountAndBySigninUser(tx *sql.Tx, g *model.GrumbleRes, signinUserId string) error {
	query := `select count(*) from bookmarks
    where grumble_pk = ?`
	// ここでtxを使うとなぜか "sql: transaction has already been committed or rolled back"
	row := s.db.QueryRow(query, g.Pk)
	row.Scan(&g.BookmarkedCount)

	query = `select count(*) from bookmarks
    where grumble_pk = ? and by_user_id = ?`
	// ここでtxを使うとなぜか "sql: transaction has already been committed or rolled back"
	row = s.db.QueryRow(query, g.Pk, signinUserId)
	count := 0
	row.Scan(&count)
	if count > 0 {
		g.IsBookmarkedBySigninUser = true
	}
	return nil
}

func (s *grumbleStore) RetrieveByUserId(signinUserId string, userId string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	res := make([]model.GrumbleRes, 0)
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    where u.id = ?`
	rows, err := tx.Query(query, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		g := model.GrumbleRes{}
		err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, g)
	}

	err = tx.Commit()
	return res, err
}

func (s *grumbleStore) Search(signinUserId string, searchWord string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	res := make([]model.GrumbleRes, 0)
	pattern := "%" + searchWord + "%"
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    where g.content like ?`
	rows, err := tx.Query(query, pattern)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		g := model.GrumbleRes{}
		err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, g)
	}

	err = tx.Commit()
	return res, err
}

func (s *grumbleStore) DeleteByPk(pk string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete grumbles where pk = ?", pk)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *grumbleStore) CreateBookmark(grumblePk string, byUserId string) (model.Bookmark, error) {
	res := model.Bookmark{
		GrumblePk: grumblePk,
		ByUserId:  byUserId,
	}
	tx, err := s.db.Begin()
	if err != nil {
		return res, err
	}
	query := `select count(*) from bookmarks
    where grumble_pk = ? and by_user_id = ?`
	row := tx.QueryRow(query, grumblePk, byUserId)
	count := 0
	row.Scan(&count)
	if count > 0 {
		tx.Rollback()
		return res, errors.New("すでにブックマークしています。")
	}

	query = `insert into bookmarks
    (grumble_pk, by_user_id, created_at)
    values (?, ?, ?)`
	_, err = tx.Exec(query, grumblePk, byUserId, time.Now())
	if err != nil {
		tx.Rollback()
		return res, err
	}

	return res, tx.Commit()
}

func (s *grumbleStore) DeleteBookmark(grumblePk string, byUserId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `delete from bookmarks
    where grumble_pk = ? and by_user_id = ?`
	_, err = tx.Exec(query, grumblePk, byUserId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// トランザクションは呼び出し元でRollBack or Commitさせる
func (s *grumbleStore) retrieveBookmarksByUserId(tx *sql.Tx, userId string) ([]model.Bookmark, error) {
	res := make([]model.Bookmark, 0)
	query := `select pk, grumble_pk, by_user_id from bookmarks
    where by_user_id = ?`
	rows, err := tx.Query(query, userId)
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

func (s *grumbleStore) RetrieveBookmarkedGrumblesByUserId(signinUserId string, userId string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	bookmarks, err := s.retrieveBookmarksByUserId(tx, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := make([]model.GrumbleRes, 0)
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    where g.pk = ?`
	for _, b := range bookmarks {
		rows, err := tx.Query(query, b.GrumblePk)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		defer rows.Close()

		// ループは一周しかしない（はず）
		for rows.Next() {
			g := model.GrumbleRes{}
			err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			res = append(res, g)
		}
	}

	return res, tx.Commit()
}
