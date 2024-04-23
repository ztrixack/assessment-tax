package database

import "database/sql"

type Database interface {
	Config() config
	Close() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryOne(query string, args ...interface{}) (*sql.Row, error)
	Execute(query string, args ...interface{}) (sql.Result, error)
}

func queryOne(db *sql.DB, query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	return row, nil
}

func query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	return rows, err
}

func execute(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(args...)
	return ret, err
}
