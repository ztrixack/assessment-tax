package logger

import (
	"io"
	"os"
)

type config struct {
	level  int
	writer io.Writer
}

func Config() *config {
	return &config{
		level:  getLevelFromEnv(),
		writer: os.Stdout,
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
