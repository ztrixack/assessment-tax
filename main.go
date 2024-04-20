package main

import (
	"log"

	"github.com/ztrixack/assessment-tax/internal/api"
	"github.com/ztrixack/assessment-tax/internal/database"
)

func main() {
	_, err := database.NewPostgresDB(database.Config())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	server := api.NewEchoAPI(api.Config())
	server.Listen()
}
