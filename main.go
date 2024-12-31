package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/lotusMind/meditation/api"
	db "github.com/lotusMind/meditation/db/sqlc"
	"github.com/lotusMind/meditation/util"
)

func main() {
	//Load configuration, use "." because main.go is in the same directory
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load configuration")
	}

	// create a new database connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not start the server", err)
	}
}
