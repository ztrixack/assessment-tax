package swagger

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
)

type handler struct {
}

func New(e api.API) {
	server := &handler{}
	server.setupRoutes(e.GetRouter())
}

func (h handler) setupRoutes(r api.Router) {
	r.GET("/swagger/*", echoSwagger.WrapHandler)
}
