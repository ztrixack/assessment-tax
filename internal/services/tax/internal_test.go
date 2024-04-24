package tax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateStepTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		lower       float64
		upper       float64
		rate        float64
		expectedTax float64
	}{
		{"Zero rate", 100000, 0, 150000, 0, 0},
		{"Within first bracket", 200000, 150000, 500000, 0.10, 5000},
		{"Within second bracket", 750000, 500000, 1000000, 0.15, 37500},
		{"Within third bracket", 1500000, 1000000, 2000000, 0.20, 100000},
		{"Within fourth bracket", 3000000, 2000000, 3000000, 0.35, 350000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax := calculateStepTax(tt.income, tt.lower, tt.upper, tt.rate)
			assert.Equal(t, tt.expectedTax, gotTax)
		})
	}
}

func TestCalculateProgressiveTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		expectedTax float64
		expectedErr error
	}{
		{
			name:        "Story: EXP01",
			income:      440000,
			expectedTax: 29000,
			expectedErr: nil,
		},
		{
			name:        "Negative income",
			income:      -100,
			expectedTax: 0,
			expectedErr: ErrNegativeIncome,
		},
		{
			name:        "Zero income",
			income:      0,
			expectedTax: 0,
			expectedErr: nil,
		},
		{
			name:        "Income within first bracket (0% tax)",
			income:      150000,
			expectedTax: 0,
			expectedErr: nil,
		},
		{
			name:        "Income within second bracket on the lower end (10% tax)",
			income:      150001,
			expectedTax: 0.1,
			expectedErr: nil,
		},
		{
			name:        "Income within second bracket on the upper end (10% tax)",
			income:      500000,
			expectedTax: 35000,
			expectedErr: nil,
		},
		{
			name:        "Income within third bracket on the lower end (15% tax)",
			income:      500001,
			expectedTax: 35000.15,
			expectedErr: nil,
		},
		{
			name:        "Income within third bracket on the upper end (15% tax)",
			income:      1000000,
			expectedTax: 110000,
			expectedErr: nil,
		},
		{
			name:        "Income within fourth bracket on the lower end (20% tax)",
			income:      1000001,
			expectedTax: 110000.2,
			expectedErr: nil,
		},
		{
			name:        "Income within fourth bracket on the upper end (20% tax)",
			income:      2000000,
			expectedTax: 310000,
			expectedErr: nil,
		},
		{
			name:        "Income within fifth bracket on the lower end (35% tax)",
			income:      2000001,
			expectedTax: 310000.35,
			expectedErr: nil,
		},
		{
			name:        "Income within fifth bracket (35% tax)",
			income:      5000000,
			expectedTax: 1360000,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax, gotErr := calculateProgressiveTax(tt.income)
			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedTax, gotTax)
		})
	}
}
