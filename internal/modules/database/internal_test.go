package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

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
			db, mock := setup(t)
			defer db.Close()

			tt.mockBehaviour(mock)
			row, err := queryOne(db, tt.query, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var id int
			err = row.Scan(&id)
			if err != sql.ErrNoRows {
				assert.NoError(t, err)
			}

			mock.ExpectationsWereMet()
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
			db, mock := setup(t)
			defer db.Close()

			tt.mockBehaviour(mock)
			_, err := query(db, tt.query, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			mock.ExpectationsWereMet()
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
			db, mock := setup(t)
			defer db.Close()

			tt.mockSetup(mock)

			result, err := execute(db, tt.query, tt.args...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			count, err := result.RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, count)

			mock.ExpectationsWereMet()
		})
	}
}
