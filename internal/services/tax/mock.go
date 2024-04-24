package tax

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Servicer = (*MockService)(nil)

type MockService struct {
	mock.Mock
}

func (m *MockService) Calculate(ctx context.Context, req CalculateRequest) *CalculateResponse {
	args := m.Called(ctx, req)
	return args.Get(0).(*CalculateResponse)
}
