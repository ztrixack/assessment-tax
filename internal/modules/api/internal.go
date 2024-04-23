package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type API interface {
	Config() config
	Listen() error
	Close() error
	GetRouter() Router
	Use(middleware ...echo.MiddlewareFunc)
	NewContext(*http.Request, http.ResponseWriter) echo.Context
}

type Router interface {
	GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	POST(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	PUT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	PATCH(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Context = echo.Context

type CustomValidator struct {
	validator *validator.Validate
}

// Validate is a method that applies the validation rules to the input struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
