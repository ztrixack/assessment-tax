package main

import (
	"github.com/ztrixack/assessment-tax/internal/domain/system"
	"github.com/ztrixack/assessment-tax/internal/infra/api"
	"github.com/ztrixack/assessment-tax/internal/infra/database"
	"github.com/ztrixack/assessment-tax/internal/infra/logger"
)

func main() {
	log := logger.NewZerolog(logger.Config())

	_, err := database.NewPostgresDB(database.Config())
	if err != nil {
		log.Err(err).C("Failed to connect to database")
	}

	server := api.NewEchoAPI(api.Config())

	system.New(server)

	server.Listen()
}
