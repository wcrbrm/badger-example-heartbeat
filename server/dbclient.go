package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database Client structure
type DbClient struct {
	Dsn            string
	DB             *sqlx.DB
}

// Connect - return a connection to a postgres database wi
func NewDatabaseClient() *DbClient {
	DSN, ok := os.LookupEnv("PGSQL_DSN")
	if !ok {
		DSN = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", DSN)
	if err != nil {
		log.Fatal("[db] Fatal Error, cannot open postgres database ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("[db] Ping Error", err)
	}
	return &DbClient{DSN, db}
}
