package tax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

func TestCalculations(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		request      CalculationsRequest
		expected     CalculationsResponse
		expectedCode int
	}{
		{
			name:        "Story: EXP01",
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
			name:        "Successful calculation",
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
			name:         "Request parameters are invalid on Bind",
			contentType:  "text/plain",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "Request parameters are invalid on Validate",
			contentType: "application/json",
			request: CalculationsRequest{
				TotalIncome: 500000.0,
				WHT:         -1000,
				Allowances:  []Allowance{},
			},
			expectedCode: http.StatusBadRequest,
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
			h := New(log, server)

			err = h.Calculations(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var result CalculationsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
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
				WHT:         0.0,
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
