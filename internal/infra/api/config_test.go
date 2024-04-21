package api

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	expectedPort := "8080"
	os.Setenv("PORT", expectedPort)
	defer os.Unsetenv("PORT")

	c := Config()
	if c.port != expectedPort {
		t.Errorf("expected '%s' but got '%s'", expectedPort, c.port)
	}
}
