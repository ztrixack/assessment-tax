package tax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

func TestToResponse(t *testing.T) {
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
