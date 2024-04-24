package admin

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

func TestDeductionsPersonalRequest(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		request      DeductionsPersonalRequest
		expected     DeductionsPersonalResponse
		expectedCode int
	}{
		{
			name:        "Story: EXP05",
			contentType: "application/json",
			request: DeductionsPersonalRequest{
				Amount: 70000.0,
			},
			expected: DeductionsPersonalResponse{
				PersonalDeduction: 70000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Normal case",
			contentType: "application/json",
			request: DeductionsPersonalRequest{
				Amount: 50000.0,
			},
			expected: DeductionsPersonalResponse{
				PersonalDeduction: 50000.0,
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
			request: DeductionsPersonalRequest{
				Amount: -70000.0,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := api.NewEchoAPI(api.Config())

			reqb, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewBuffer(reqb))
			req.Header.Set("Content-Type", tt.contentType)
			req.Header.Set("Accept", tt.contentType)
			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)
			res := rec.Result()
			defer res.Body.Close()

			log := logger.NewMockLogger()
			h := New(log, server)

			err = h.DeductionsPersonal(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var result DeductionsPersonalResponse
				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeductionsPersonalRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request DeductionsPersonalRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			request: DeductionsPersonalRequest{
				Amount: 70000.0,
			},
			wantErr: false,
		},
		{
			name: "Amount on the lower end",
			request: DeductionsPersonalRequest{
				Amount: 10000.0,
			},
			wantErr: false,
		},
		{
			name: "Amount on the upper end",
			request: DeductionsPersonalRequest{
				Amount: 100000.0,
			},
			wantErr: false,
		},
		{
			name: "Amount on the below lower end",
			request: DeductionsPersonalRequest{
				Amount: 9999.0,
			},
			wantErr: true,
		},
		{
			name: "Amount on the above upper end",
			request: DeductionsPersonalRequest{
				Amount: 100001.0,
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
