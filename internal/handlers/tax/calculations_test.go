package tax

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

func TestCalculations(t *testing.T) {
	tests := []struct {
		name         string
		mockBehavior func(*tax.MockService)
		contentType  string
		request      CalculationsRequest
		expected     CalculationsResponse
		expectedCode int
	}{
		{
			name: "Story: EXP01",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 29000.0}, nil)
			},
			contentType: "application/json",
			request: CalculationsRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax: 29000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Story: EXP02",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 4000.0}, nil)
			},
			contentType: "application/json",
			request: CalculationsRequest{
				TotalIncome: 500000.0,
				WHT:         25000.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax: 4000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Successful calculation",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 90000.0}, nil)
			},
			contentType: "application/json",
			request: CalculationsRequest{
				TotalIncome: 1500000.0,
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax: 90000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Request parameters are invalid on Bind",
			mockBehavior: func(ms *tax.MockService) {},
			contentType:  "text/plain",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Request parameters are invalid on Validate",
			mockBehavior: func(ms *tax.MockService) {},
			contentType:  "application/json",
			request: CalculationsRequest{
				TotalIncome: -500000.0,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Calculation service is broken",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))
			},
			contentType: "application/json",
			request: CalculationsRequest{
				TotalIncome: 1500000.0,
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := api.NewEchoAPI(api.Config())

			reqb, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqb))
			req.Header.Set("Content-Type", tt.contentType)
			req.Header.Set("Accept", tt.contentType)
			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)
			res := rec.Result()
			defer res.Body.Close()

			log := logger.NewMockLogger()
			ms := new(tax.MockService)
			h := New(log, server, ms)

			tt.mockBehavior(ms)
			err = h.Calculations(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var result CalculationsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			ms.AssertExpectations(t)
		})
	}
}

func TestCalculationRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CalculationsRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CalculationsRequest{
				TotalIncome: 500000.0,
				WHT:         50000.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0}},
			},
			wantErr: false,
		},
		{
			name: "no wht",
			request: CalculationsRequest{
				TotalIncome: 500000.0,
			},
			wantErr: false,
		},
		{
			name: "no allowances",
			request: CalculationsRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
			},
			wantErr: false,
		},
		{
			name: "invalid total income",
			request: CalculationsRequest{
				TotalIncome: -1,
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			err := v.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
