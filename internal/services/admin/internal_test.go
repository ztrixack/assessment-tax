package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		request  SetDeductionRequest
		wantErr  bool
		errValue error
	}{
		{"Valid Personal Request", SetDeductionRequest{Type: Personal, Amount: 50000}, false, nil},
		{"Invalid Personal Request", SetDeductionRequest{Type: Personal, Amount: 1000}, true, ErrLessThanLimit(Personal, 10000)},
		{"Valid KReceipt Request", SetDeductionRequest{Type: KReceipt, Amount: 5000}, false, nil},
		{"Invalid KReceipt Request", SetDeductionRequest{Type: KReceipt, Amount: 1000000}, true, ErrMoreThanLimit(KReceipt, 100000)},
		{"Unknown Type", SetDeductionRequest{Type: DeductionType("Unknown"), Amount: 10000}, true, ErrInvalidDeductionType},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.validate()
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.errValue, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLimiter(t *testing.T) {
	tests := []struct {
		name    string
		dtype   DeductionType
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
			err := limiter(tc.dtype, tc.amount, tc.lower, tc.upper)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeType(t *testing.T) {
	tests := []struct {
		name      string
		dtype     DeductionType
		wantValue string
	}{
		{"Map personal", Personal, "personal"},
		{"Map k-receipt", KReceipt, "k_receipt"},
		{"Map k-receipt as string", "k-receipt", "k_receipt"},
		{"Map anything else", "unknown", "unknown"},
		{"Map anything else with hyphen", "un-known", "un_known"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantValue, sanitizeType(tc.dtype))
		})
	}
}
