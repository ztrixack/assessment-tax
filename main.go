package main

import (
	"github.com/ztrixack/assessment-tax/internal/handlers/admin"
	"github.com/ztrixack/assessment-tax/internal/handlers/swagger"
	"github.com/ztrixack/assessment-tax/internal/handlers/system"
	"github.com/ztrixack/assessment-tax/internal/handlers/tax"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/api/middlewares"
	"github.com/ztrixack/assessment-tax/internal/modules/database"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	admin_service "github.com/ztrixack/assessment-tax/internal/services/admin"
	tax_service "github.com/ztrixack/assessment-tax/internal/services/tax"

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
	// modules
	log := logger.NewZerolog(logger.Config())
	server := api.NewEchoAPI(api.Config())
	server.Use(middlewares.Logger(log))
	db, err := database.NewPostgresDB(database.Config())
	if err != nil {
		log.Err(err).C("Failed to connect to database")
	}
	defer db.Close()

	// services
	taxService := tax_service.New(log, db)
	adminService := admin_service.New(log, db)

	// handlers
	system.New(server)
	swagger.New(server)
	tax.New(log, server, taxService)
	admin.New(log, server, adminService)

	// application
	log.Fields(logger.Fields{"port": server.Config().Port}).I("Starting server")
	if err := server.Listen(); err != nil {
		log.Err(err).C("Failed to start server")
	}
	log.I("Stopping server")
}
