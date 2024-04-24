package admin

import (
	"context"
	"fmt"
	"unicode"
)

type Servicer interface {
	SetDeduction(ctx context.Context, request SetDeductionRequest) (float64, error)
}

type DeductionType string

const (
	Personal DeductionType = "personal"
	KReceipt DeductionType = "k-receipt"

	PersonalMinimum = 10000
	PersonalMaximum = 100000

	KReceiptMinimum = 0
	KReceiptMaximum = 100000
)

var (
	ErrInvalidDeductionType = fmt.Errorf("invalid deduction type")
	ErrLessThanLimit        = func(dtype DeductionType, value float64) error {
		return fmt.Errorf("the %s deduction cannot be less than %f", dtype, value)
	}
	ErrMoreThanLimit = func(dtype DeductionType, value float64) error {
		return fmt.Errorf("the %s deduction cannot be more than %f", dtype, value)
	}
	ErrUpdateDatabase = func(dtype DeductionType) error {
		return fmt.Errorf("failed to set %s deduction", dtype)
	}
)

func (r SetDeductionRequest) validate() error {
	switch r.Type {
	case Personal:
		return limiter(Personal, r.Amount, PersonalMinimum, PersonalMaximum)

	case KReceipt:
		return limiter(KReceipt, r.Amount, KReceiptMinimum, KReceiptMaximum)
	}

	return ErrInvalidDeductionType
}

func limiter(dtype DeductionType, amount, lower, upper float64) error {
	if amount < lower {
		return ErrLessThanLimit(dtype, lower)
	}

	if amount > upper {
		return ErrMoreThanLimit(dtype, upper)
	}

	return nil
}

func sanitizeType(dtype DeductionType) string {
	sanitized := make([]rune, len(dtype))
	for i, r := range dtype {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sanitized[i] = r
		} else {
			sanitized[i] = '_'
		}
	}

	return string(sanitized)
}
