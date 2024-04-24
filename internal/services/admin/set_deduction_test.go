package admin

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ztrixack/assessment-tax/internal/modules/database"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

func setup() (*service, sqlmock.Sqlmock, error) {
	log := logger.NewMockLogger()
	db, mock, err := database.NewMockDB()
	return &service{log, db}, mock, err
}

func TestSetDeduction(t *testing.T) {
	s, mock, err := setup()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer s.db.Close()

	tests := []struct {
		name          string
		request       SetDeductionRequest
		mockBehaviour func()
		expectedError error
	}{
		{
			name:    "Successful to set personal update",
			request: SetDeductionRequest{Type: Personal, Amount: 60000.0},
			mockBehaviour: func() {
				mock.ExpectPrepare("UPDATE allowances SET personal = \\$1").
					ExpectExec().
					WithArgs(60000.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:    "Successful to set k-receipt update",
			request: SetDeductionRequest{Type: KReceipt, Amount: 60000.0},
			mockBehaviour: func() {
				mock.ExpectPrepare("UPDATE allowances SET k_receipt = \\$1").
					ExpectExec().
					WithArgs(60000.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:          "Set Personal deduction less than 10,000",
			request:       SetDeductionRequest{Type: Personal, Amount: 9999.0},
			mockBehaviour: func() {},
			expectedError: ErrLessThanLimit(Personal, PersonalMinimum),
		},
		{
			name:          "Set Personal deduction more than 100,000",
			request:       SetDeductionRequest{Type: Personal, Amount: 100001.0},
			mockBehaviour: func() {},
			expectedError: ErrMoreThanLimit(Personal, PersonalMaximum),
		},
		{
			name:          "Set K-Receipt deduction less than 0",
			request:       SetDeductionRequest{Type: KReceipt, Amount: -1.0},
			mockBehaviour: func() {},
			expectedError: ErrLessThanLimit(KReceipt, KReceiptMinimum),
		},
		{
			name:          "Set K-Receipt deduction more than 100,000",
			request:       SetDeductionRequest{Type: KReceipt, Amount: 100001.0},
			mockBehaviour: func() {},
			expectedError: ErrMoreThanLimit(KReceipt, KReceiptMaximum),
		},
		{
			name:    "Database error",
			request: SetDeductionRequest{Type: Personal, Amount: 60000.0},
			mockBehaviour: func() {
				mock.ExpectPrepare("UPDATE allowances SET personal = \\$1").
					ExpectExec().
					WithArgs(60000.0).
					WillReturnError(assert.AnError)
			},
			expectedError: ErrUpdateDatabase(Personal),
		},
		{
			name:          "Unknown Type",
			request:       SetDeductionRequest{Type: "unknown", Amount: 60000.0},
			mockBehaviour: func() {},
			expectedError: ErrInvalidDeductionType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			_, err := s.SetDeduction(context.Background(), tt.request)

			assert.Equal(t, tt.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
