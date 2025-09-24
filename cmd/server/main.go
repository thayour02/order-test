package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sava/cmd/api"
	db "github.com/sava/db/sqlc"
	"github.com/sava/env"
)

var (
	address = ":8080"
	dbDriver = "postgres"
	DbSource string
)

func init() {
	DbSource = env.Getenv("DB_SOURCE", "dbsource")

	// fallback for local development
	if DbSource == "" {
		DbSource = "postgresql://root:secret@localhost:5434/order?sslmode=disable"
	}
}

func main() {
	conn, err := sql.Open(dbDriver, DbSource)
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