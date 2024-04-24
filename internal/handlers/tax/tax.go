package tax

import (
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

type handler struct {
	log logger.Logger
}

func New(log logger.Logger, e api.API) *handler {
	handler := &handler{log}
	handler.setupRoutes(e.GetRouter())
	return handler
}

func (h handler) setupRoutes(r api.Router) {
	r.POST("/tax/calculations", h.Calculations)
}
