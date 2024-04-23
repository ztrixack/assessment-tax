package database

import "os"

type config struct {
	DatabaseURL string
}

func Config() *config {
	return &config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
