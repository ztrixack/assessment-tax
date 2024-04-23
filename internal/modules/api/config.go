package api

import "os"

type config struct {
	Host string
	Port string
}

func Config() *config {
	return &config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
}
