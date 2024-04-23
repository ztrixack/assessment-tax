package main

import (
	"github.com/ztrixack/assessment-tax/internal/handlers/swagger"
	"github.com/ztrixack/assessment-tax/internal/handlers/system"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/database"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"

	_ "github.com/ztrixack/assessment-tax/docs"
)

// @title			Assessment Tax API
// @version		1.0
// @description	Assessment Tax API for Go Bootcamp
//
// @contact.name	Tanawat Hongthai
// @contact.url	https://github.com/ztrixack/assessment-tax.git
// @contact.email	ztrixack.th@gmail.com
//
// @schemes		http
func main() {
	log := logger.NewZerolog(logger.Config())

	_, err := database.NewPostgresDB(database.Config())
	if err != nil {
		log.Err(err).C("Failed to connect to database")
	}

	server := api.NewEchoAPI(api.Config())

	system.New(server)
	swagger.New(server)

	log.I("Starting server")
	server.Listen()
	log.I("Stopping server")
}
