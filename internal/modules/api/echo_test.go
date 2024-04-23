package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func setup() *echoAPI {
	c := &config{Host: "localhost", Port: "9999"}
	return NewEchoAPI(c)
}

func TestNewEchoAPI(t *testing.T) {
	api := setup()
	defer api.Close()

	if api.config.Port != "9999" {
		t.Errorf("Expected API port to be '9999' but got '%s'", api.config.Port)
	}
	if api.router == nil {
		t.Errorf("Router should not be nil")
	}
}

func TestEchoAPIRouting(t *testing.T) {
	api := setup()
	defer api.Close()

	// Define a test handler
	api.router.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test passed")
	})

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Serve HTTP request using the Echo router
	api.router.ServeHTTP(rec, req)

	// Check the status code and response
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}
	if rec.Body.String() != "test passed" {
		t.Errorf("Unexpected body %q", rec.Body.String())
	}
}

func TestEchoAPI_Notify(t *testing.T) {
	server := setup()
	defer server.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Listen(); err != nil {
			t.Errorf("Expected no error on server listen but got '%v'", err)
		}
	}()

	time.Sleep(1 * time.Second)
	stop <- os.Interrupt
	time.Sleep(1 * time.Second)

	if err := server.Close(); err != nil {
		t.Errorf("Expected no error on server close but got '%v'", err)
	}
}
