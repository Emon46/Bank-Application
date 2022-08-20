package main

import (
	"database/sql"
	"github.com/emon46/bank-application/api"
	db "github.com/emon46/bank-application/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable"
	ServerAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
