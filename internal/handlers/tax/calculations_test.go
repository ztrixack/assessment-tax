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
	"github.com/ztrixack/assessment-tax/internal/utils/constants"
)

func pointerTo(value float64) *float64 {
	return &value
}

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
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 29000.0, Refund: 0.0, TaxLevel: []float64{0.0, 29000.0, 0.0, 0.0, 0.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax:      29000.0,
				TaxLevel: []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 29000.0}, {constants.T500k_1M, 0.0}, {constants.T1M_2M, 0.0}, {constants.T2M, 0.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Story: EXP02",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 4000.0, Refund: 0.0, TaxLevel: []float64{0.0, 29000.0, 0.0, 0.0, 0.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         25000.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax:      4000.0,
				TaxLevel: []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 29000.0}, {constants.T500k_1M, 0.0}, {constants.T1M_2M, 0.0}, {constants.T2M, 0.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Story: EXP03",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 19000.0, Refund: 0.0, TaxLevel: []float64{0.0, 19000.0, 0.0, 0.0, 0.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 200000.0}},
			},
			expected: CalculationsResponse{
				Tax:      19000.0,
				TaxLevel: []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 19000.0}, {constants.T500k_1M, 0.0}, {constants.T1M_2M, 0.0}, {constants.T2M, 0.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Story: EXP04",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 19000.0, Refund: 0.0, TaxLevel: []float64{0.0, 19000.0, 0.0, 0.0, 0.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 200000.0}},
			},
			expected: CalculationsResponse{
				Tax:      19000.0,
				TaxLevel: []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 19000.0}, {constants.T500k_1M, 0.0}, {constants.T1M_2M, 0.0}, {constants.T2M, 0.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Story: EXP07",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 14000.0, Refund: 0.0, TaxLevel: []float64{0.0, 14000.0, 0.0, 0.0, 0.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "k-receipt", Amount: 200000.0}, {AllowanceType: "donation", Amount: 100000.0}},
			},
			expected: CalculationsResponse{
				Tax:      14000.0,
				TaxLevel: []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 14000.0}, {constants.T500k_1M, 0.0}, {constants.T1M_2M, 0.0}, {constants.T2M, 0.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Successful calculation",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 1360000.0, Refund: 0.0, TaxLevel: []float64{0.0, 35000.0, 75000.0, 200000.0, 1050000.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(5000000.0),
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax:      1360000.0,
				TaxLevel: []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 35000.0}, {constants.T500k_1M, 75000.0}, {constants.T1M_2M, 200000.0}, {constants.T2M, 1050000.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Successful with Refund",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 0.0, Refund: 21000.0, TaxLevel: []float64{0.0, 29000.0, 0.0, 0.0, 0.0}}, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         50000.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			},
			expected: CalculationsResponse{
				Tax:       0.0,
				TaxRefund: pointerTo(21000.0),
				TaxLevel:  []TaxLevel{{constants.T0_150k, 0.0}, {constants.T150k_500k, 29000.0}, {constants.T500k_1M, 0.0}, {constants.T1M_2M, 0.0}, {constants.T2M, 0.0}},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Request parameters are invalid on Bind",
			mockBehavior: func(ms *tax.MockService) {
				// Do nothing
			},
			contentType:  constants.TEXT_PLAIN,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Request parameters are invalid on Validate",
			mockBehavior: func(ms *tax.MockService) {
				// Do nothing
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(-500000.0),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Calculation service is broken",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))
			},
			contentType: constants.APPLICATION_JSON,
			request: CalculationsRequest{
				TotalIncome: pointerTo(1500000.0),
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

func TestCalculationRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CalculationsRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         50000.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 200000.0}},
			},
			wantErr: false,
		},
		{
			name: "no wht",
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
			},
			wantErr: false,
		},
		{
			name: "no allowances",
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         0.0,
			},
			wantErr: false,
		},
		{
			name: "invalid total income",
			request: CalculationsRequest{
				TotalIncome: pointerTo(-1),
				WHT:         0.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: 0}},
			},
			wantErr: true,
		},
		{
			name: "invalid allowance type",
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         5000.0,
				Allowances:  []Allowance{{AllowanceType: "unknown", Amount: 10000}},
			},
			wantErr: true,
		},
		{
			name: "invalid allowance amount",
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         5000.0,
				Allowances:  []Allowance{{AllowanceType: "donation", Amount: -500}},
			},
			wantErr: true,
		},
		{
			name:    "No request",
			request: CalculationsRequest{},
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

func TestCalculationsRequestToServiceRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  CalculationsRequest
		expected tax.CalculateRequest
	}{
		{
			name: "valid request",
			request: CalculationsRequest{
				TotalIncome: pointerTo(500000.0),
				WHT:         5000.0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 10000},
				},
			},
			expected: tax.CalculateRequest{
				Income: 500000.0,
				WHT:    5000.0,
				Allowances: []tax.Allowance{
					{Type: tax.AllowanceType("donation"), Amount: 10000},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.toServiceRequest()
			assert.Equal(t, tt.expected, result)
		})
	}
}
