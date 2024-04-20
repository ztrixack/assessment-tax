package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type echoAPI struct {
	router *echo.Echo
}

func NewEchoAPI() *echoAPI {
	server := &echoAPI{
		router: echo.New(),
	}
	server.setAppHandlers()
	return server
}

func (s *echoAPI) Listen() error {
	return s.router.Start(":1323")
}

func (s *echoAPI) Close() {
	s.router.Shutdown(context.Background())
}

func (s *echoAPI) setAppHandlers() {
	s.router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
}
