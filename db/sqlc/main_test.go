package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"github.com/sava/env"
	_ "github.com/lib/pq"
)


var testQueries *Queries
var (
	dbDriver = "postgres"
	Dbsource string
)

func init() {
	Dbsource :=  env.Getenv("DB_SOURCE", "dbsource")
	if Dbsource == "" {
		Dbsource = "postgresql://root:secret@localhost:5434/order?sslmode=disable" // fallback for local
	}
}
func TestMain(m *testing.M) {
	com, err := sql.Open(dbDriver, Dbsource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(com)

	os.Exit(m.Run())
}