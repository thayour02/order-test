package main


import (
	"database/sql"
	"log"
	"github.com/sava/env"

	_ "github.com/lib/pq"
	"github.com/sava/cmd/api"
	db "github.com/sava/db/sqlc"
)

	
	var (
		address = ":8080"
	dbDriver = "postgres"
	Dbsource string

)

func init() {
	Dbsource :=  env.Getenv("DB_SOURCE", "dbsource")
	if Dbsource == "" {
		Dbsource = "postgresql://root:secret@localhost:5434/order?sslmode=disable" // fallback for local
	}
}

func main() {
	conn, err := sql.Open(dbDriver, Dbsource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}