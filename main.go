package main

import (
	"GoBankProject/api"
	db "GoBankProject/db/sqlc"
	"GoBankProject/util"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)
// TODO: Refactor code to load all conf from env variables or conf file

func main() {
	config, err := util.LoadConfig("../GoBankProject"); if err != nil{
		log.Fatal("cannot load configurations", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource);if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(conn)
	sever := api.NewServer(store)

	err = sever.Run(config.ServerAddress); if err != nil{
		log.Fatal("cannot start server", err)
	}
}
