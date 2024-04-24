package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var _ Database = (*postgresDB)(nil)

type postgresDB struct {
	config *config
	db     *sql.DB
}

func NewPostgresDB(c *config) (*postgresDB, error) {
	db, err := sql.Open("postgres", c.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresDB{c, db}, nil
}

func (p *postgresDB) Config() config {
	return *p.config
}

func (p *postgresDB) Close() error {
	return p.db.Close()
}

func (p *postgresDB) Query(qry string, args ...interface{}) (*sql.Rows, error) {
	return query(p.db, qry, args...)
}

func (p *postgresDB) QueryOne(query string, args ...interface{}) (*sql.Row, error) {
	return queryOne(p.db, query, args...)
}

func (p *postgresDB) Execute(query string, args ...interface{}) (sql.Result, error) {
	return execute(p.db, query, args...)
}
