package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/admin"
	"github.com/ztrixack/assessment-tax/internal/utils/constants"
)

func TestDeductionsPersonalRequest(t *testing.T) {
	tests := []struct {
		name         string
		mockBehavior func(*admin.MockService)
		contentType  string
		request      DeductionsPersonalRequest
		expected     DeductionsPersonalResponse
		expectedCode int
	}{
		{
			name: "Story: EXP05",
			mockBehavior: func(ms *admin.MockService) {
				ms.On("SetDeduction", mock.Anything, mock.Anything).Return(70000.0, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: DeductionsPersonalRequest{
				Amount: 70000.0,
			},
			expected: DeductionsPersonalResponse{
				PersonalDeduction: 70000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Normal case",
			mockBehavior: func(ms *admin.MockService) {
				ms.On("SetDeduction", mock.Anything, mock.Anything).Return(50000.0, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: DeductionsPersonalRequest{
				Amount: 50000.0,
			},
			expected: DeductionsPersonalResponse{
				PersonalDeduction: 50000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Request parameters are invalid on Bind",
			mockBehavior: func(ms *admin.MockService) {
				// Do nothing
			},
			contentType:  constants.TEXT_PLAIN,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Request parameters are invalid on Validate",
			mockBehavior: func(ms *admin.MockService) {
				// Do nothing
			},
			contentType: constants.APPLICATION_JSON,
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
			ms := new(admin.MockService)
			h := New(log, server, ms)

			tt.mockBehavior(ms)
			err = h.DeductionsPersonal(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var result DeductionsPersonalResponse
				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			ms.AssertExpectations(t)
		})
	}
}

func TestDeductionsPersonalRequestValidation(t *testing.T) {
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

func TestDeductionsPersonalRequestToServiceRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  DeductionsPersonalRequest
		expected admin.SetDeductionRequest
	}{
		{
			name: "valid request",
			request: DeductionsPersonalRequest{
				Amount: 50000.0,
			},
			expected: admin.SetDeductionRequest{
				Type:   "personal",
				Amount: 50000.0,
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
