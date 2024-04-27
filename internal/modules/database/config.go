package database

import "os"

const DEFAULT_DATABASE_URL = "host=localhost port=5432 user=porsgres password=porsgres dbname=ktaxes sslmode=disable"

type config struct {
	DatabaseURL string
}

func Config() *config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = DEFAULT_DATABASE_URL
	}

	return &config{
		DatabaseURL: databaseURL,
	}
}
