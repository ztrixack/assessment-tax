package api

import (
	"os"
	"testing"

	"github.com/go-playground/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		expectedHost string
		expectedPort string
	}{
		{
			name:         "Normal case",
			env:          map[string]string{"HOST": "", "PORT": "8080"},
			expectedHost: "",
			expectedPort: "8080",
		},
		{
			name:         "PORT set to 8082",
			env:          map[string]string{"PORT": "8082"},
			expectedHost: "localhost",
			expectedPort: "8082",
		},
		{
			name:         "non-numeric PORT",
			env:          map[string]string{"PORT": "non-numeric"},
			expectedHost: "localhost",
			expectedPort: "8080",
		},
		{
			name:         "set empty HOST",
			env:          map[string]string{"HOST": ""},
			expectedHost: "",
			expectedPort: "8080",
		},
		{
			name:         "set any HOST",
			env:          map[string]string{"HOST": "new-host"},
			expectedHost: "new-host",
			expectedPort: "8080",
		},
		{
			name:         "no ENV set",
			env:          map[string]string{},
			expectedHost: "localhost",
			expectedPort: "8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			c := Config()

			assert.Equal(t, tt.expectedHost, c.Host)
			assert.Equal(t, tt.expectedPort, c.Port)
		})
	}
}
