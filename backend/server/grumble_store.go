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
	err = s.retrieveReplyInfo(tx, &res)
	if err != nil {
		tx.Rollback()
		return res, err
	}
	err = s.retrieveBookmarkedCountAndBySigninUser(tx, &res, signinUserId)
	if err != nil {
		tx.Rollback()
		return res, err
	}
	return res, tx.Commit()
}

func (s *grumbleStore) RetrieveReplyAncestors(grumblePk string, signinUserId string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    left join replies as r
        on r.dst_grumble_pk = g.pk
    where r.src_grumble_pk = ?`
	res := make([]model.GrumbleRes, 0)

	var rec func(gpk string) error
	rec = func(gpk string) error {
		g := model.GrumbleRes{}
		row := tx.QueryRow(query, gpk)
		err = row.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
		if err != nil {
			// tx.Rollback()
			return nil
		}
		err = s.retrieveReplyInfo(tx, &g)
		if err != nil {
			// tx.Rollback()
			return err
		}
		err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			// tx.Rollback()
			return err
		}
		res = append(res, g)
		return rec(g.Pk)
	}

	err = rec(grumblePk)
	if err != nil {
		return nil, tx.Rollback()
	}
	return res, tx.Commit()
}

func (s *grumbleStore) RetrieveByReplyDstPk(grumblePk string, signinUserId string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    left join replies as r
        on r.src_grumble_pk = g.pk
    where r.dst_grumble_pk = ?`
	rows, err := tx.Query(query, grumblePk)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	res := make([]model.GrumbleRes, 0)
	for rows.Next() {
		g := model.GrumbleRes{}
		err = rows.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = s.retrieveReplyInfo(tx, &g)
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
	return res, tx.Commit()
}

func (s *grumbleStore) Create(content string, user model.User) (model.Grumble, error) {
	res := model.Grumble{}
	t := time.Now()
	id, err := createUlid(t)
	if err != nil {
		return res, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return res, err
	}
	res = model.Grumble{
		Pk:        id,
		Content:   content,
		UserId:    user.Id,
		CreatedAt: t,
	}
	_, err = tx.Exec("insert into grumbles (pk, content, user_id, created_at) values (?, ?, ?, ?)", res.Pk, res.Content, res.UserId, res.CreatedAt)
	if err != nil {
		tx.Rollback()
		return res, err
	}
	err = tx.Commit()
	return res, err
}

func (s *grumbleStore) retrieveReplyInfo(tx *sql.Tx, g *model.GrumbleRes) error {
	query := `select count(*) from replies
    where dst_grumble_pk = ?`
	row := s.db.QueryRow(query, g.Pk)
	err := row.Scan(&g.Reply.RepliedCount)
	if err != nil {
		return err
	}

	query = `select r.dst_grumble_pk, g.user_id from replies as r
    left join grumbles as g
        on r.dst_grumble_pk = g.pk
    where r.src_grumble_pk = ?`
	row = s.db.QueryRow(query, g.Pk)
	err = row.Scan(&g.Reply.DstGrumblePk, &g.Reply.DstUserId)
	return nil
}

func (s *grumbleStore) retrieveRegrumbleInfo(tx *sql.Tx, g *model.GrumbleRes, signinUserId string) error {
	query := `select count(*) from regrumbles
    where grumble_pk = ?`
	rows, err := tx.Query(query, g.Pk)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.Regrumble.RegrumbledCount)
		if err != nil {
			return err
		}
		break
	}

	query = `select `
	return nil
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

func (s *grumbleStore) retrieveRegrumbledCountAndBySigninUser(tx *sql.Tx, g *model.GrumbleRes, signinUserId string) error {
	query := `select count(*) from regrumbles
    where grumble_pk = ?`
	rows, err := s.db.Query(query, g.Pk)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.Regrumble.RegrumbledCount)
		if err != nil {
			return err
		}
		break
	}

	query = `select count(*) from regrumbles
    where grumble_pk = ? and by_user_id = ?`
	rows, err = s.db.Query(query, g.Pk, signinUserId)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			g.Regrumble.IsRegrumbledBySigninUser = true
		}
		break
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
		err = s.retrieveReplyInfo(tx, &g)
		if err != nil {
			tx.Rollback()
			return res, err
		}
		err = s.retrieveRegrumbledCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return res, err
		}
		err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, g)
	}

	// リグランブル取得
	query = `select g.pk, g.content, g.user_id, g.created_at, u.name, re.created_at
    from regrumbles as re
    left join grumbles as g
        on g.pk = re.grumble_pk
    left join users as u
        on g.user_id = u.id
    where u.id = ?`
	rows2, err := tx.Query(query, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		g := model.GrumbleRes{}
		g.Regrumble.IsRegrumble = true
		err = rows2.Scan(&g.Pk, &g.Content, &g.UserId, &g.CreatedAt, &g.UserName, &g.Regrumble.CreatedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = s.retrieveReplyInfo(tx, &g)
		if err != nil {
			tx.Rollback()
			return res, err
		}
		err = s.retrieveRegrumbledCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return res, err
		}
		err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, g)
	}
	return res, tx.Commit()
}

func (s *grumbleStore) RetrieveRegrumblesByUserId(signinUserId string, userId string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	res := make([]model.GrumbleRes, 0)

	return res, tx.Commit()
}

func (s *grumbleStore) Search(signinUserId string, searchWord string) ([]model.GrumbleRes, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	res := make([]model.GrumbleRes, 0)
	query := `select g.pk, g.content, g.user_id, g.created_at, u.name
    from grumbles as g
    left join users as u
        on g.user_id = u.id
    where g.content like concat('%', ?, '%')`
	rows, err := tx.Query(query, searchWord)
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
		err = s.retrieveReplyInfo(tx, &g)
		if err != nil {
			tx.Rollback()
			return res, err
		}
		err = s.retrieveBookmarkedCountAndBySigninUser(tx, &g, signinUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, g)
	}

	return res, tx.Commit()
}

func (s *grumbleStore) Delete(pk string, userId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `delete g, r, b from grumbles as g
    left join replies as r
        on r.src_grumble_pk = g.pk
    left join bookmarks as b
        on b.grumble_pk = g.pk
    where g.pk = ? and g.user_id = ?`
	_, err = tx.Exec(query, pk, userId)
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
			err = s.retrieveReplyInfo(tx, &g)
			if err != nil {
				tx.Rollback()
				return res, err
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

func (s *grumbleStore) CreateReply(srcGrumblePk string, dstGrumblePk string) (model.Reply, error) {
	res := model.Reply{}
	tx, err := s.db.Begin()
	if err != nil {
		return res, err
	}
	query := `insert into replies
    (created_at, src_grumble_pk, dst_grumble_pk)
    values (?, ?, ?)`
	_, err = tx.Exec(query, time.Now(), srcGrumblePk, dstGrumblePk)
	if err != nil {
		tx.Rollback()
		return res, err
	}

	return res, tx.Commit()
}

func (s *grumbleStore) CreateRegrumble(grumblePk string, byUserId string) (model.Regrumble, error) {
	res := model.Regrumble{}
	tx, err := s.db.Begin()
	if err != nil {
		return res, err
	}
	query := `select user_id from grumbles
    where pk = ?`
	row := tx.QueryRow(query, grumblePk)
	if row.Err() != nil {
		tx.Rollback()
		return res, row.Err()
	}
	var dstUserId string
	err = row.Scan(&dstUserId)
	if err != nil {
		tx.Rollback()
		return res, err
	}
	now := time.Now()
	query = `insert into regrumbles
    (created_at, grumble_pk, dst_user_id, by_user_id)
    values (?, ?, ?, ?)`
	_, err = tx.Exec(query, now, grumblePk, dstUserId, byUserId)
	if err != nil {
		tx.Rollback()
		return res, err
	}

	return res, tx.Commit()
}

func (s *grumbleStore) DeleteRegrumble(grumblePk string, byUserId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `delete from regrumbles
    where grumble_pk = ? and by_user_id = ?`
	_, err = tx.Exec(query, grumblePk, byUserId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
