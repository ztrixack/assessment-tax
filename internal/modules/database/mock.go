package database

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

var _ Database = (*sqlMockDB)(nil)

type sqlMockDB struct {
	db *sql.DB
}

func NewMockDB() (*sqlMockDB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	return &sqlMockDB{db}, mock, nil
}

func (p *sqlMockDB) Config() config {
	return config{}
}

func (p *sqlMockDB) Close() error {
	return p.db.Close()
}

func (p *sqlMockDB) Query(query_ string, args ...interface{}) (*sql.Rows, error) {
	return query(p.db, query_, args...)
}

func (p *sqlMockDB) QueryOne(query string, args ...interface{}) (*sql.Row, error) {
	return queryOne(p.db, query, args...)
}

func (p *sqlMockDB) Execute(query string, args ...interface{}) (sql.Result, error) {
	return execute(p.db, query, args...)
}
