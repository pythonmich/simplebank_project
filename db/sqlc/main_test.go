package db

import (
	"GoBankProject/util"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)


var testQueries *Queries
var testDB *sql.DB
var config util.Config

func TestMain(m *testing.M) {
	var err error
	config, err = util.LoadConfig("../.."); if err != nil {
		log.Fatal("cannot load configurations", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil{
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
