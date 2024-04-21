package domain

import (
	"github.com/ztrixack/assessment-tax/internal/infra/api"
)

type System interface {
	Root(api.Context) error
}
