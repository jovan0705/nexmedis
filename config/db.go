package config

import (
    "context"
    "fmt"
    "log"
    "github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func GetDB() (*pgx.Conn) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/nexmedis_db")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil
	}

	// Connection successful
	fmt.Println("Successfully connected to the database!")
	return conn
}
