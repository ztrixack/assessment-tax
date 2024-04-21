package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	db *sql.DB
}

func NewPostgresDB(c *config) (*postgresDB, error) {
	db, err := sql.Open("postgres", c.database_url)
	if err != nil {
		return nil, err
	}

	return &postgresDB{db}, nil
}
