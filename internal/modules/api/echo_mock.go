package api

import (
	"io"
	"net/http"

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

func (s *echoMockAPI) Config() config {
	return config{}
}

func (s *echoMockAPI) Listen() error {
	panic("unimplemented")
}

func (s *echoMockAPI) Close() error {
	panic("unimplemented")
}

func (s *echoMockAPI) GetRouter() Router {
	return s.router
}

func (s *echoMockAPI) Use(middleware ...echo.MiddlewareFunc) {
	s.router.Use(middleware...)
}
func (s *echoMockAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *echoMockAPI) NewContext(r *http.Request, w http.ResponseWriter) echo.Context {
	return s.router.NewContext(r, w)
}
