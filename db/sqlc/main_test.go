package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sava/env"
)


var testQueries *Queries
// 	const (dbDriver = "postgres")

// 	var dbSource string

// func init () {
// 	dbSource = env.Getenv("DB_SOURCE", "dbsource")
// 	if dbSource == "" {
// 		dbSource = "postgresql://root:secret@localhost:5434/order?sslmode=disable" // fallback for local
// 	}
// }

var dbDriver = "postgres"
var dbSource = func() string {
    if src := env.Getenv("DB_SOURCE", ""); src != "" {
        return src
    }
    return "postgresql://root:secret@localhost:5434/order?sslmode=disable"
}()


func TestMain(m *testing.M) {
	com, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(com)

	os.Exit(m.Run())
}