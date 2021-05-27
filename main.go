package main

import (
	"GoBankProject/api"
	db "GoBankProject/db/sqlc"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)
// TODO: Refactor code to load all conf from env variables or conf file
const (
	dbDriver = "postgres"
	dbSource = "postgresql://python_mich:Musyimi7.@localhost:5432/simple_bank?sslmode=disable"
	severAddress = "0.0.0.0:5050"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource);if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(conn)
	sever := api.NewServer(store)

	err = sever.Run(severAddress); if err != nil{
		log.Fatal("cannot start server", err)
	}
}
