package tax

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name           string
		request        CalculateRequest
		expectedResult *CalculateResponse
	}{
		{
			name: "Story: EXP01",
			request: CalculateRequest{
				Income:     500000.0,
				Allowances: []Allowance{},
			},
			expectedResult: &CalculateResponse{
				Tax: 29000.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewMockLogger()
			s := New(log)

			ctx := context.Background()

			result := s.Calculate(ctx, tt.request)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
