package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func CreateConnection() (*sql.DB, error) {
	port := 5432
	host := "postgres"

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, "postgres", "flights", "postgres")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
