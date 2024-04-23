package logger

import (
	"os"
	"testing"
)

func TestGetLevelFromEnv(t *testing.T) {
	tests := []struct {
		envValue string
		expected int
	}{
		{"debug", 0},
		{"info", 1},
		{"warn", 2},
		{"error", 3},
		{"invalid", 1}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.envValue, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.envValue)
			defer os.Unsetenv("LOG_LEVEL")

			level := getLevelFromEnv()
			if level != tt.expected {
				t.Errorf("expected '%d' but got '%d'", tt.expected, level)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	defer os.Unsetenv("LOG_LEVEL")

	config := Config()
	if config.Level != 1 {
		t.Errorf("expected '1' but got '%d'", config.Level)
	}
	if config.Writer != os.Stdout {
		t.Errorf("expected 'os.Stdout' but got '%v'", config.Writer)
	}
}
