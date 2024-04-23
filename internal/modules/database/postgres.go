package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var _ Database = (*postgresDB)(nil)

type postgresDB struct {
	db *sql.DB
}

func NewPostgresDB(c *config) (*postgresDB, error) {
	db, err := sql.Open("postgres", c.database_url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresDB{db}, nil
}
