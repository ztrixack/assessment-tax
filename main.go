package main

import (
	"github.com/ztrixack/assessment-tax/internal/api"
	"github.com/ztrixack/assessment-tax/internal/database"
	"github.com/ztrixack/assessment-tax/internal/logger"
)

func main() {
	log := logger.NewZerolog(logger.Config())

	_, err := database.NewPostgresDB(database.Config())
	if err != nil {
		log.Err(err).C("Failed to connect to database")
	}

	server := api.NewEchoAPI(api.Config())
	server.Listen()

}
