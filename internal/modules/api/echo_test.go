package api

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func setup() *echoAPI {
	config := &config{port: "9999"}
	server := NewEchoAPI(config)
	return server
}

func TestNewEchoAPI(t *testing.T) {
	port := "8080"
	server := NewEchoAPI(&config{port})

	if server.config.port != port {
		t.Errorf("Expected API port to be '%s' but got '%s'", port, server.config.port)
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

func TestEchoAPI_Shutdown(t *testing.T) {
	server := setup()

	go func() {
		server.Listen()
	}()

	time.Sleep(1 * time.Second)

	// Shutdown the server
	err := server.Close()
	if err != nil {
		t.Errorf("Expected no error on server shutdown but got '%v'", err)
	}
}
