package admin

import (
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/api/middlewares"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/admin"
)

type handler struct {
	log   logger.Logger
	admin admin.Servicer
}

func New(log logger.Logger, e api.API, admin admin.Servicer) *handler {
	handler := &handler{log, admin}
	handler.setupRoutes(e.GetRouter())
	return handler
}

func (h handler) setupRoutes(r api.Router) {
	r.POST("/admin/deductions/personal", h.DeductionsPersonal, middlewares.BasicAuth(h.log))
	r.POST("/admin/deductions/k-receipt", h.DeductionsKReceipt, middlewares.BasicAuth(h.log))
}
