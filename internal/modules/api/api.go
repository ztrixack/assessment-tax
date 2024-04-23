package api

import "github.com/labstack/echo/v4"

type API interface {
	Listen() error
	Close() error
	GetRouter() *echo.Echo
}

type Router interface {
	GET(string, echo.HandlerFunc)
}

type Context = echo.Context
