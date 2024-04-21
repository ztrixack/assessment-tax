package api

import (
	"github.com/labstack/echo/v4"
	"github.com/ztrixack/assessment-tax/internal/domain/system"
)

var _ system.API = (*api)(nil)

type api struct {
}

func New(e *echo.Echo) system.API {
	server := &api{}
	server.setupRoutes(e)
	return server
}

func (s api) setupRoutes(r *echo.Echo) {
	r.GET("/", s.Root)
}
