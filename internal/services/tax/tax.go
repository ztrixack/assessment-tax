package tax

import (
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

var _ Servicer = (*service)(nil)

type service struct {
	log logger.Logger
}

func New(log logger.Logger) *service {
	services := &service{log}
	return services
}
