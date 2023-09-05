package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mhg14/simplebank/api"
	db "github.com/mhg14/simplebank/db/sqlc"
	"github.com/mhg14/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server:", err)
	}
	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("can not start server", err)
	}
}
