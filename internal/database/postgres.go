package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	db *sql.DB
}

func NewPostgresDB() (*postgresDB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &postgresDB{db}, nil
}
