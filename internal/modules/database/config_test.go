package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		expectedUrl string
	}{
		{
			name: "default",
			env: map[string]string{
				"DATABASE_URL": "host=porsgres port=5432 user=porsgres password=porsgres dbname=porsgres sslmode=disable",
			},
			expectedUrl: "host=porsgres port=5432 user=porsgres password=porsgres dbname=porsgres sslmode=disable",
		},
		{
			name:        "no URL set",
			env:         map[string]string{},
			expectedUrl: "host=localhost port=5432 user=porsgres password=porsgres dbname=ktaxes sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			c := Config()

			assert.Equal(t, tt.expectedUrl, c.DatabaseURL)
		})
	}
}
