//go:build unit

package database

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	expectedURL := "postgres://username:password@localhost:5432/dbname"
	os.Setenv("DATABASE_URL", expectedURL)
	defer os.Unsetenv("DATABASE_URL")

	c := Config()
	if c.database_url != expectedURL {
		t.Errorf("expected '%s' but got '%s'", expectedURL, c.database_url)
	}
}
