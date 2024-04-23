package database

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	expectedURL := "host=porsgres port=5432 user=porsgres password=porsgres dbname=porsgres sslmode=disable"
	os.Setenv("DATABASE_URL", expectedURL)
	defer os.Unsetenv("DATABASE_URL")

	c := Config()
	if c.DatabaseURL != expectedURL {
		t.Errorf("expected '%s' but got '%s'", expectedURL, c.DatabaseURL)
	}
}
