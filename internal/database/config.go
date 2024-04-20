package database

import "os"

type config struct {
	database_url string
}

func Config() *config {
	return &config{
		database_url: os.Getenv("DATABASE_URL"),
	}
}
