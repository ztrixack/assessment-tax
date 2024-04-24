package tax

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

func TestUploadCSV(t *testing.T) {
	tests := []struct {
		name         string
		mockBehavior func(*tax.MockService)
		taxFile      *multipart.FileHeader
		expected     UploadCSVResponse
		expectedCode int
	}{
		{
			name: "Story: EXP06",
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 29000.0}, nil).Once()
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 25000.0}, nil).Once()
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 0.0}, nil).Once()
			},
			expectedCode: http.StatusOK,
			expected: UploadCSVResponse{
				Taxes: []Tax{
					{TotalIncome: 500000, Tax: 29000},
					{TotalIncome: 600000, Tax: 25000},
					{TotalIncome: 750000, Tax: 0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := api.NewEchoAPI(api.Config())

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			fileField, _ := writer.CreateFormFile("taxFile", "taxes.csv")
			fileField.Write([]byte("totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000"))
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)
			res := rec.Result()
			defer res.Body.Close()

			log := logger.NewMockLogger()
			ms := new(tax.MockService)
			h := New(log, server, ms)

			tt.mockBehavior(ms)
			err := h.UploadCSV(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var result UploadCSVResponse
				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			ms.AssertExpectations(t)
		})
	}
}
