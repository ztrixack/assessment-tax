package system

import (
	"github.com/ztrixack/assessment-tax/internal/domain"
	"github.com/ztrixack/assessment-tax/internal/infra/api"
)

var _ domain.System = (*systemAPI)(nil)

type systemAPI struct {
}

func New(e api.API) domain.System {
	server := &systemAPI{}
	server.setupRoutes(e)
	return server
}

func (s systemAPI) setupRoutes(r api.API) {
	r.GetRouter().GET("/", s.Root)
}
