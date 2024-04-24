package tax

import (
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

type handler struct {
	log logger.Logger
	tax tax.Servicer
}

func New(log logger.Logger, e api.API, tax tax.Servicer) *handler {
	handler := &handler{log, tax}
	handler.setupRoutes(e.GetRouter())
	return handler
}

func (h handler) setupRoutes(r api.Router) {
	r.POST("/tax/calculations", h.Calculations)
	r.POST("/tax/calculations/upload-csv", h.UploadCSV)
}
