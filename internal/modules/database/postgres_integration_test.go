//go:build integration

package database

import (
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		hasError bool
	}{
		{"Database connection successful", "host=postgres-tests port=5432 user=test password=test dbname=testdb sslmode=disable", false},
		{"Database connection failed", "host=invalid port=5432 user=test password=test dbname=testdb sslmode=disable", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPostgresDB(&config{database_url: tt.url})
			if tt.hasError != (err != nil) {
				t.Errorf("expected error but got '%v'", err)
			}
		})
	}
}
