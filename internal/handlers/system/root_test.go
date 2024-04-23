package system

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
)

func TestRootHandler(t *testing.T) {
	e := api.NewEchoMockAPI()
	New(e)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.GetRouter().ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "Hello, Go Bootcamp!"
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}
