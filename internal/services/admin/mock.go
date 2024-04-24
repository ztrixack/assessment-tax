package admin

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Servicer = (*MockService)(nil)

type MockService struct {
	mock.Mock
}

func (m *MockService) SetDeduction(ctx context.Context, req SetDeductionRequest) (float64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(float64), args.Error(1)
}
