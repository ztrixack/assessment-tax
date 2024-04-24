package tax

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	defaultMockBehavior := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"personal"}).AddRow(60000)
		mock.ExpectPrepare("SELECT personal FROM allowances").ExpectQuery().WillReturnRows(rows)
	}

	tests := []struct {
		name           string
		mockBehavior   func(mock sqlmock.Sqlmock)
		request        CalculateRequest
		expectedResult *CalculateResponse
		wantErr        bool
	}{
		{
			name: "Story: EXP01",
			request: CalculateRequest{
				Income:     500000.0,
				Allowances: []Allowance{},
			},
			mockBehavior: defaultMockBehavior,
			expectedResult: &CalculateResponse{
				Tax: 29000.0,
			},
			wantErr: false,
		},
		{
			name: "error in tax calculation",
			request: CalculateRequest{
				Income:     -1,
				Allowances: []Allowance{},
			},
			mockBehavior:   defaultMockBehavior,
			expectedResult: nil,
			wantErr:        true,
		},
		{
			name: "error in database",
			request: CalculateRequest{
				Income:     500000.0,
				Allowances: []Allowance{},
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT personal FROM allowances").ExpectQuery().WillReturnError(errors.New("some error"))
			},
			expectedResult: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svr, mock, close := setup(t)
			defer close()

			tt.mockBehavior(mock)
			result, err := svr.Calculate(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
