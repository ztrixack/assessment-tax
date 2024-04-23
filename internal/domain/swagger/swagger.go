package swagger

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/ztrixack/assessment-tax/internal/infra/api"
)

type handler struct {
}

func New(e api.API) {
	server := &handler{}
	server.setupRoutes(e)
}

func (h handler) setupRoutes(r api.API) {
	r.GetRouter().GET("/swagger/*", echoSwagger.WrapHandler)
}
