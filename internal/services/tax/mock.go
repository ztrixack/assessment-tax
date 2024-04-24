package tax

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Servicer = (*MockService)(nil)

type MockService struct {
	mock.Mock
}

func (m *MockService) Calculate(ctx context.Context, req CalculateRequest) (*CalculateResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CalculateResponse), args.Error(1)
}
