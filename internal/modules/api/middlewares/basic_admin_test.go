package middlewares

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCredentials(t *testing.T) {
	os.Setenv("ADMIN_USERNAME", "adminTax")
	os.Setenv("ADMIN_PASSWORD", "admin!")
	defer func() {
		os.Unsetenv("ADMIN_USERNAME")
		os.Unsetenv("ADMIN_PASSWORD")
	}()

	tests := []struct {
		name     string
		username string
		password string
		want     bool
	}{
		{"valid credentials", "adminTax", "admin!", true},
		{"invalid username", "invalidUser", "admin!", false},
		{"invalid password", "adminTax", "invalidPassword", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCredentials(tt.username, tt.password)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSecurePassword(t *testing.T) {
	tests := []struct {
		name   string
		pwd    string
		result string
	}{
		{"short password", "12345", "******"},
		{"long password", "password123", "pass****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			masked := securePassword(tt.pwd)
			assert.Equal(t, tt.result, masked)
		})
	}
}
