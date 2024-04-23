package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var _ API = (*echoAPI)(nil)

type echoAPI struct {
	config *config
	router *echo.Echo
}

func NewEchoAPI(c *config) *echoAPI {
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Validator = &CustomValidator{validator: validator.New()}
	server := &echoAPI{
		config: c,
		router: e,
	}
	return server
}

func (s *echoAPI) Config() config {
	return *s.config
}

func (s *echoAPI) Listen() error {
	var reterr error
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
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

func (s *echoAPI) GetRouter() Router {
	return s.router
}

func (s *echoAPI) Use(middleware ...echo.MiddlewareFunc) {
	s.router.Use(middleware...)
}

func (s *echoAPI) NewContext(r *http.Request, w http.ResponseWriter) echo.Context {
	return s.router.NewContext(r, w)
}
