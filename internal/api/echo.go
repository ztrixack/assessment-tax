package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type echoAPI struct {
	config *config
	router *echo.Echo
}

func NewEchoAPI(c *config) *echoAPI {
	server := &echoAPI{
		config: c,
		router: echo.New(),
	}
	server.setAppHandlers()
	return server
}

func (s *echoAPI) Listen() error {
	port := fmt.Sprintf(":%s", s.config.port)
	return s.router.Start(port)
}

func (s *echoAPI) Close() {
	s.router.Shutdown(context.Background())
}

func (s *echoAPI) setAppHandlers() {
	s.router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
}
