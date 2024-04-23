package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestQueryOne(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		args          []interface{}
		mockBehaviour func(mock sqlmock.Sqlmock)
		wantErr       bool
	}{
		{
			name:  "Successful query",
			query: "SELECT id FROM users WHERE username = ?",
			args:  []interface{}{"john"},
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectPrepare("SELECT id FROM users WHERE username = \\?").
					ExpectQuery().WithArgs("john").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "Query with preparation error",
			query: "SELECT id FROM users WHERE username = ?",
			args:  []interface{}{"john"},
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT id FROM users WHERE username = \\?").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockBehaviour(mock)

			row, err := queryOne(db, tt.query, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("queryOne() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				var id int
				err = row.Scan(&id)
				if err != nil && err != sql.ErrNoRows {
					t.Errorf("Failed to scan row: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestQuery(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		args          []interface{}
		mockBehaviour func(mock sqlmock.Sqlmock)
		wantErr       bool
	}{
		{
			name:  "Successful query multiple rows",
			query: "SELECT * FROM users WHERE age > ?",
			args:  []interface{}{25},
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "age"}).
					AddRow(1, "john", 30).
					AddRow(2, "jane", 27)
				mock.ExpectPrepare("SELECT \\* FROM users WHERE age > \\?").
					ExpectQuery().WithArgs(25).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "Query with preparation error",
			query: "SELECT * FROM users",
			args:  nil,
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT \\* FROM users").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name:  "Query execution error",
			query: "SELECT * FROM users WHERE username = ?",
			args:  []interface{}{"nonexistent"},
			mockBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT \\* FROM users WHERE username = \\?").
					ExpectQuery().WithArgs("nonexistent").WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockBehaviour(mock)

			_, err = query(db, tt.query, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("query() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	var tests = []struct {
		name      string
		query     string
		args      []interface{}
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
		wantCount int64
	}{
		{
			name:  "Execute insert successfully",
			query: "INSERT INTO users(username, age) VALUES (?, ?)",
			args:  []interface{}{"johndoe", 30},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO users\\(username, age\\) VALUES \\(\\?, \\?\\)").
					ExpectExec().WithArgs("johndoe", 30).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:  "Execute update with no rows affected",
			query: "UPDATE users SET age = ? WHERE username = ?",
			args:  []interface{}{31, "janedoe"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("UPDATE users SET age = \\? WHERE username = \\?").
					ExpectExec().WithArgs(31, "janedoe").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:  "Execution failure due to SQL error",
			query: "DELETE FROM users WHERE username = ?",
			args:  []interface{}{"nonexistent"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("DELETE FROM users WHERE username = \\?").
					ExpectExec().WithArgs("nonexistent").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name:  "Preparation failure due to SQL error",
			query: "DELETE FROM users WHERE username = ?",
			args:  nil,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("DELETE FROM users WHERE username = \\?").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			result, err := execute(db, tt.query, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				count, err := result.RowsAffected()
				if err != nil {
					t.Errorf("error getting rows affected: %v", err)
				}
				if count != tt.wantCount {
					t.Errorf("expected %d affected rows, got %d", tt.wantCount, count)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
