package tax

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

func TestToCalculationsResponse(t *testing.T) {
	pointerTo := func(value float64) *float64 {
		return &value
	}

	tests := []struct {
		name     string
		input    tax.CalculateResponse
		expected CalculationsResponse
	}{
		{
			name: "no refund",
			input: tax.CalculateResponse{
				Tax:      100.0,
				Refund:   0.0,
				TaxLevel: []float64{0, 0, 0, 0, 0},
			},
			expected: CalculationsResponse{
				Tax:      100.0,
				TaxLevel: []TaxLevel{{"0-150,000", 0}, {"150,001-500,000", 0}, {"500,001-1,000,000", 0}, {"1,000,001-2,000,000", 0}, {"2,000,001 ขึ้นไป", 0}},
			},
		},
		{
			name: "with refund",
			input: tax.CalculateResponse{
				Tax:      100.0,
				Refund:   50.0,
				TaxLevel: []float64{0, 0, 0, 0, 0},
			},
			expected: CalculationsResponse{
				Tax:       100.0,
				TaxRefund: pointerTo(50.0),
				TaxLevel:  []TaxLevel{{"0-150,000", 0}, {"150,001-500,000", 0}, {"500,001-1,000,000", 0}, {"1,000,001-2,000,000", 0}, {"2,000,001 ขึ้นไป", 0}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := toCalculationsResponse(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRemapAllowances(t *testing.T) {
	tests := []struct {
		name     string
		input    []Allowance
		expected []tax.Allowance
	}{
		{
			name:     "empty input",
			input:    []Allowance{},
			expected: []tax.Allowance{},
		},
		{
			name: "single element",
			input: []Allowance{
				{AllowanceType: "donation", Amount: 1000},
			},
			expected: []tax.Allowance{
				{Type: tax.AllowanceType("donation"), Amount: 1000},
			},
		},
		{
			name: "multiple elements",
			input: []Allowance{
				{AllowanceType: "donation", Amount: 1000},
				{AllowanceType: "donation", Amount: 100000},
				{AllowanceType: "unknown", Amount: 300},
			},
			expected: []tax.Allowance{
				{Type: tax.AllowanceType("donation"), Amount: 1000},
				{Type: tax.AllowanceType("donation"), Amount: 100000},
				{Type: tax.AllowanceType("unknown"), Amount: 300},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := remapAllowances(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemapTaxLevels(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		levels   []float64
		expected []TaxLevel
	}{
		{
			name:     "matching lengths",
			labels:   []string{"0-150,000", "150,001-500,000", "500,001-1,000,000"},
			levels:   []float64{10.0, 20.0, 30.0},
			expected: []TaxLevel{{"0-150,000", 10.0}, {"150,001-500,000", 20.0}, {"500,001-1,000,000", 30.0}},
		},
		{
			name:     "labels fewer than levels",
			labels:   []string{"0-150,000", "150,001-500,000"},
			levels:   []float64{10.0, 20.0, 30.0},
			expected: []TaxLevel{{"Bucket: #1", 10.0}, {"Bucket: #2", 20.0}, {"Bucket: #3", 30.0}},
		},
		{
			name:     "labels more than levels",
			labels:   []string{"0-150,000", "150,001-500,000", "500,001-1,000,000", "1,000,001-2,000,000"},
			levels:   []float64{10.0, 20.0, 30.0},
			expected: []TaxLevel{{"Bucket: #1", 10.0}, {"Bucket: #2", 20.0}, {"Bucket: #3", 30.0}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := remapTaxLevel(tc.labels, tc.levels)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetFileFromRequest(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name    string
		setup   func() *echo.Context
		wantErr bool
	}{
		{
			name: "With file",
			setup: func() *echo.Context {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				fileField, _ := writer.CreateFormFile("taxFile", "dummy.txt")
				fileField.Write([]byte("dummy data"))
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/", body)
				req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return &c
			},
			wantErr: false,
		},
		{
			name: "Without file",
			setup: func() *echo.Context {
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return &c
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			file, err := getFileFromRequest(*c)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, file)
			}
		})
	}
}

func TestParseCSVFile(t *testing.T) {
	tests := []struct {
		name           string
		file           *multipart.FileHeader
		expectedResult []tax.CalculateRequest
		wantErr        bool
	}{
		{
			name: "Dummy Test",
			file: &multipart.FileHeader{},
			expectedResult: []tax.CalculateRequest{
				{
					Income:     500000.0,
					WHT:        0.0,
					Allowances: []tax.Allowance{{Type: "donation", Amount: 0.0}},
				},
				{
					Income:     600000.0,
					WHT:        40000.0,
					Allowances: []tax.Allowance{{Type: "donation", Amount: 20000.0}},
				},
				{
					Income:     750000.0,
					WHT:        50000.0,
					Allowances: []tax.Allowance{{Type: "donation", Amount: 15000.0}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCSVFile(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestCalculateTaxes(t *testing.T) {
	pointerTo := func(value float64) *float64 {
		return &value
	}

	tests := []struct {
		name          string
		requests      []tax.CalculateRequest
		mockBehavior  func(*tax.MockService)
		expectedTaxes []Tax
		wantErr       bool
	}{
		{
			name: "successful calculations",
			requests: []tax.CalculateRequest{
				{Income: 500000.0},
				{Income: 600000.0, WHT: 40000.0},
				{Income: 500000.0, WHT: 50000.0, Allowances: []tax.Allowance{{Type: tax.Donation, Amount: 15000.0}}},
			},
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 29000.0}, nil).Once()
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 25000.0}, nil).Once()
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 0.0}, nil).Once()
			},
			expectedTaxes: []Tax{
				{TotalIncome: 500000.0, Tax: 29000.0},
				{TotalIncome: 600000.0, Tax: 25000.0},
				{TotalIncome: 500000.0, Tax: 0.0},
			},
			wantErr: false,
		},
		{
			name: "successful with refund",
			requests: []tax.CalculateRequest{
				{Income: 600000.0, WHT: 100000.0},
			},
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 0.0, Refund: 50000.0}, nil).Once()
			},
			expectedTaxes: []Tax{
				{TotalIncome: 600000.0, Tax: 0.0, TaxRefund: pointerTo(50000.0)},
			},
			wantErr: false,
		},
		{
			name: "error on second calculation",
			requests: []tax.CalculateRequest{
				{Income: 500000.0},
				{Income: 600000.0, WHT: 40000.0},
			},
			mockBehavior: func(ms *tax.MockService) {
				ms.On("Calculate", mock.Anything, mock.Anything).Return(&tax.CalculateResponse{Tax: 29000.0}, nil).Once()
				ms.On("Calculate", mock.Anything, mock.Anything).Return(nil, errors.New("calculation error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			server := api.NewEchoAPI(api.Config())
			log := logger.NewMockLogger()
			ms := new(tax.MockService)
			h := New(log, server, ms)

			tc.mockBehavior(ms)
			gotTaxes, err := h.calculateTaxes(ctx, tc.requests)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTaxes, gotTaxes)
			}

			ms.AssertExpectations(t)
		})
	}
}

func TestToTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		response    tax.CalculateResponse
		expectedTax Tax
	}{
		{
			name:   "Successfully case",
			income: 500000.0,
			response: tax.CalculateResponse{
				Tax:    29000.0,
				Refund: 0,
			},
			expectedTax: Tax{
				TotalIncome: 500000.0,
				Tax:         29000.0,
			},
		},
		{
			name:   "Successfully case with refund",
			income: 500000.0,
			response: tax.CalculateResponse{
				Tax:    0.0,
				Refund: 1000.0,
			},
			expectedTax: Tax{
				TotalIncome: 500000.0,
				Tax:         0.0,
				TaxRefund:   &[]float64{1000.0}[0],
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := toTax(tc.income, tc.response)
			assert.Equal(t, tc.expectedTax, result)
		})
	}
}
