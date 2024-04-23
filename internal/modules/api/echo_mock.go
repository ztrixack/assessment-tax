package api

import (
	"io"

	"github.com/labstack/echo/v4"
)

var _ API = (*echoMockAPI)(nil)

type echoMockAPI struct {
	router *echo.Echo
}

func NewEchoMockAPI() *echoMockAPI {
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	return &echoMockAPI{
		router: e,
	}
}

func (s *echoMockAPI) Listen() error {
	panic("unimplemented")
}

func (s *echoMockAPI) Close() error {
	panic("unimplemented")
}

func (s *echoMockAPI) GetRouter() *echo.Echo {
	return s.router
}
