package main


import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sava/cmd/api"
	db "github.com/sava/db/sqlc"
)

const (
	address = ":8080"
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5434/order?sslmode=disable"
)

func main() {
conn, err := sql.Open(dbDriver, dbSource)
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