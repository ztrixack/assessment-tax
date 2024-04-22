package system

import (
	"net/http"

	"github.com/ztrixack/assessment-tax/internal/infra/api"
)

// Root godoc
//
//	@summary		Hello, Go Bootcamp!
//	@description	Hello, Go Bootcamp!
//	@tags			system
//	@produce		text/plain
//	@success		200	{string}	string	"Hello, Go Bootcamp!"
//	@router			/ [get]
func (a systemAPI) Root(c api.Context) error {
	return c.String(http.StatusOK, "Hello, Go Bootcamp!")
}
