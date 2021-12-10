package main

import (
	"log"

	"github.com/guricerin/grumbler/backend/db"
	"github.com/guricerin/grumbler/backend/server"
	"github.com/guricerin/grumbler/backend/util"
)

func main() {
	cfg, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.OpenMySql(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(cfg, db)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
