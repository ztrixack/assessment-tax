package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	var reterr error
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := fmt.Sprintf(":%s", s.config.port)
		if err := s.router.Start(port); err != nil && err != http.ErrServerClosed {
			reterr = err
			stop <- os.Interrupt
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.router.Shutdown(ctx); err != nil {
		return err
	}

	return reterr
}

func (s *echoAPI) Close() error {
	return s.router.Shutdown(context.Background())
}

func (s *echoAPI) setAppHandlers() {
	s.router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
}
