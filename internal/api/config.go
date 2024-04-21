package api

import "os"

type config struct {
	port string
}

func Config() *config {
	return &config{
		port: os.Getenv("PORT"),
	}
}
