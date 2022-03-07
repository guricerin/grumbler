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
	query := `insert into follows
    (src_user_id, dst_user_id)
    values (?, ?)`
	_, err := s.db.Exec(query, src, dst)
	return err
}

func (s *followStore) Delete(src string, dst string) error {
	query := `delete from follows
    where src_user_id = ? and dst_user_id = ?`
	_, err := s.db.Exec(query, src, dst)
	return err
}

func (s *followStore) RetrieveFollows(srcUserId string) ([]model.Follow, error) {
	query := `select pk, src_user_id, dst_user_id from follows
    where src_user_id = ?`
	rows, err := s.db.Query(query, srcUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]model.Follow, 0)
	for rows.Next() {
		f := model.Follow{}
		err := rows.Scan(&f.Pk, &f.SrcUserId, &f.DstUserId)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (s *followStore) RetrieveFollowers(dstUserId string) ([]model.Follow, error) {
	query := `select pk, src_user_id, dst_user_id from follows
    where dst_user_id = ?`
	rows, err := s.db.Query(query, dstUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]model.Follow, 0)
	for rows.Next() {
		f := model.Follow{}
		err := rows.Scan(&f.Pk, &f.SrcUserId, &f.DstUserId)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

// userId1からみて、userId2はフォローなのかフォロワーなのか
func (s *followStore) RetrieveFollowRelation(userId1 string, userId2 string) (bool, bool, error) {
	var count int
	query := `select count(*) from follows
    where src_user_id = ? and dst_user_id = ?`
	row := s.db.QueryRow(query, userId1, userId2)
	err := row.Scan(&count)
	if err != nil {
		return false, false, err
	}
	isFollow := false
	if count > 0 {
		isFollow = true
	}

	row = s.db.QueryRow(query, userId2, userId1)
	err = row.Scan(&count)
	if err != nil {
		return false, false, err
	}
	isFollower := false
	if count > 0 {
		isFollower = true
	}

	return isFollow, isFollower, nil
}
