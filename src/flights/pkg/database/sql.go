package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// image: library/postgres:13
// container_name: postgres
// restart: on-failure
// environment:
//   POSTGRES_USER: postgres
//   POSTGRES_PASSWORD: "postgres"
//   POSTGRES_DB: postgres
// volumes:
//   - db-data:/var/lib/postgresql/data
//   - ./postgres/:/docker-entrypoint-initdb.d/
// ports:

func CreateConnection() (*sql.DB, error) {
	port := 5432
	host := "localhost"

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
