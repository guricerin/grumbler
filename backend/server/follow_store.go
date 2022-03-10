package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guricerin/grumbler/backend/model"
)

type followStore struct {
	db *sql.DB
}

func NewFollowStore(db *sql.DB) followStore {
	return followStore{db: db}
}

func (s *followStore) Create(src string, dst string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `insert into follows
    (src_user_id, dst_user_id)
    values (?, ?)`
	_, err = tx.Exec(query, src, dst)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *followStore) Delete(src string, dst string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `delete from follows
    where src_user_id = ? and dst_user_id = ?`
	_, err = tx.Exec(query, src, dst)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *followStore) RetrieveFollows(srcUserId string) ([]model.Follow, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	query := `select pk, src_user_id, dst_user_id from follows
    where src_user_id = ?`
	rows, err := tx.Query(query, srcUserId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	res := make([]model.Follow, 0)
	for rows.Next() {
		f := model.Follow{}
		err := rows.Scan(&f.Pk, &f.SrcUserId, &f.DstUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, f)
	}
	return res, tx.Commit()
}

func (s *followStore) RetrieveFollowers(dstUserId string) ([]model.Follow, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	query := `select pk, src_user_id, dst_user_id from follows
    where dst_user_id = ?`
	rows, err := tx.Query(query, dstUserId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	res := make([]model.Follow, 0)
	for rows.Next() {
		f := model.Follow{}
		err := rows.Scan(&f.Pk, &f.SrcUserId, &f.DstUserId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = append(res, f)
	}
	return res, tx.Commit()
}

// userId1からみて、userId2はフォローなのかフォロワーなのか
func (s *followStore) RetrieveFollowRelation(userId1 string, userId2 string) (bool, bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, false, err
	}
	var count int
	query := `select count(*) from follows
    where src_user_id = ? and dst_user_id = ?`
	row := tx.QueryRow(query, userId1, userId2)
	err = row.Scan(&count)
	if err != nil {
		tx.Rollback()
		return false, false, err
	}
	isFollow := false
	if count > 0 {
		isFollow = true
	}

	row = tx.QueryRow(query, userId2, userId1)
	err = row.Scan(&count)
	if err != nil {
		tx.Rollback()
		return false, false, err
	}
	isFollower := false
	if count > 0 {
		isFollower = true
	}

	return isFollow, isFollower, tx.Commit()
}
