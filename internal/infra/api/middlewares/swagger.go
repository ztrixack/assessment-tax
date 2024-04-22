package middlewares

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Swagger() echo.HandlerFunc {
	return echoSwagger.WrapHandler
}
