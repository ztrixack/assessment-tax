package system

import (
	"github.com/ztrixack/assessment-tax/internal/modules/api"
)

type handler struct {
}

func New(e api.API) *handler {
	server := &handler{}
	server.setupRoutes(e)
	return server
}

func (s handler) setupRoutes(r api.API) {
	r.GetRouter().GET("/", s.Root)
}
