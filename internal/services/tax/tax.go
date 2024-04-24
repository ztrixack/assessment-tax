package tax

import (
	"github.com/ztrixack/assessment-tax/internal/modules/database"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

var _ Servicer = (*service)(nil)

type service struct {
	log logger.Logger
	db  database.Database
}

func New(log logger.Logger, db database.Database) *service {
	services := &service{log, db}
	return services
}
