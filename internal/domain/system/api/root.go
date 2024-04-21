package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a api) Root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Go Bootcamp!")
}
