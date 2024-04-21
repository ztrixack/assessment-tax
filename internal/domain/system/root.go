package system

import (
	"net/http"

	"github.com/ztrixack/assessment-tax/internal/infra/api"
)

func (a systemAPI) Root(c api.Context) error {
	return c.String(http.StatusOK, "Hello, Go Bootcamp!")
}
