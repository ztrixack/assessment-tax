package logger

import (
	"io"
	"os"
)

type config struct {
	Level  int
	Writer io.Writer
}

func Config() *config {
	return &config{
		Level:  getLevelFromEnv(),
		Writer: os.Stdout,
	}
}

func getLevelFromEnv() int {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return 0
	case "info":
		return 1
	case "warn":
		return 2
	case "error":
		return 3
	default:
		return 1
	}
}
