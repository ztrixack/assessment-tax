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
				Tax:    100.0,
				Refund: 0.0,
			},
			expected: CalculationsResponse{
				Tax: 100.0,
			},
		},
		{
			name: "with refund",
			input: tax.CalculateResponse{
				Tax:    100.0,
				Refund: 50.0,
			},
			expected: CalculationsResponse{
				Tax:       100.0,
				TaxRefund: pointerTo(50.0),
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
