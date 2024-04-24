package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) {
	tests := []struct {
		name    string
		type_   DeductionType
		amount  float64
		lower   float64
		upper   float64
		wantErr bool
	}{
		{"Valid Personal Deduction", Personal, 50000, PersonalMinimum, PersonalMaximum, false},
		{"Too Low Personal Deduction", Personal, 9999, PersonalMinimum, PersonalMaximum, true},
		{"Too High Personal Deduction", Personal, 100001, PersonalMinimum, PersonalMaximum, true},
		{"Valid K-Receipt Deduction", KReceipt, 50000, KReceiptMinimum, KReceiptMaximum, false},
		{"Too Low K-Receipt Deduction", KReceipt, -1, KReceiptMinimum, KReceiptMaximum, true},
		{"Too High K-Receipt Deduction", KReceipt, 100001, KReceiptMinimum, KReceiptMaximum, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := limiter(tc.type_, tc.amount, tc.lower, tc.upper)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
