package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://python_mich:Musyimi7.@localhost:5432/simple_bank?sslmode=disable"
)
var testQueries *Queries
var testDB *sql.DB
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil{
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
