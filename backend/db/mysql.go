package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guricerin/grumbler/backend/util"
)

func OpenMySql(cfg util.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DbUrl)
	if err != nil {
		return nil, err
	}
	if err := tryPing(db); err != nil {
		return nil, err
	}

	return db, nil
}

func tryPing(conn *sql.DB) (err error) {
	for i := 0; i < 3; i++ {
		err = conn.Ping()
		if err == nil {
			return
		}
		time.Sleep(5 * time.Second)
	}
	return
}
