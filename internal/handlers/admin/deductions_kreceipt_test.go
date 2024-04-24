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

func TestDeductionsKReceiptRequest(t *testing.T) {
	tests := []struct {
		name         string
		mockBehavior func(*admin.MockService)
		contentType  string
		request      DeductionsKReceiptRequest
		expected     DeductionsKReceiptResponse
		expectedCode int
	}{
		{
			name: "Story: EXP08",
			mockBehavior: func(ms *admin.MockService) {
				ms.On("SetDeduction", mock.Anything, mock.Anything).Return(70000.0, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: DeductionsKReceiptRequest{
				Amount: 70000.0,
			},
			expected: DeductionsKReceiptResponse{
				KReceipt: 70000.0,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Normal case",
			mockBehavior: func(ms *admin.MockService) {
				ms.On("SetDeduction", mock.Anything, mock.Anything).Return(50000.0, nil)
			},
			contentType: constants.APPLICATION_JSON,
			request: DeductionsKReceiptRequest{
				Amount: 50000.0,
			},
			expected: DeductionsKReceiptResponse{
				KReceipt: 50000.0,
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
			request: DeductionsKReceiptRequest{
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

			req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewBuffer(reqb))
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
			err = h.DeductionsKReceipt(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var result DeductionsKReceiptResponse
				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			ms.AssertExpectations(t)
		})
	}
}

func TestDeductionsKReceiptRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request DeductionsKReceiptRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			request: DeductionsKReceiptRequest{
				Amount: 70000.0,
			},
			wantErr: false,
		},
		{
			name: "Amount on the lower end",
			request: DeductionsKReceiptRequest{
				Amount: 1.0,
			},
			wantErr: false,
		},
		{
			name: "Amount on the upper end",
			request: DeductionsKReceiptRequest{
				Amount: 100000.0,
			},
			wantErr: false,
		},
		{
			name: "Amount on the below lower end",
			request: DeductionsKReceiptRequest{
				Amount: 0.0,
			},
			wantErr: true,
		},
		{
			name: "Amount on the above upper end",
			request: DeductionsKReceiptRequest{
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

func TestDeductionsKReceiptRequestToServiceRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  DeductionsKReceiptRequest
		expected admin.SetDeductionRequest
	}{
		{
			name: "valid request",
			request: DeductionsKReceiptRequest{
				Amount: 50000.0,
			},
			expected: admin.SetDeductionRequest{
				Type:   "k-receipt",
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
