package system

import (
	"github.com/labstack/echo/v4"
)

type API interface {
	Root(echo.Context) error
}
