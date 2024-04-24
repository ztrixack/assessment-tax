package admin

import (
	"context"
	"fmt"
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
	ErrLessThanLimit        = func(type_ DeductionType, value float64) error {
		return fmt.Errorf("the %s deduction cannot be less than %f", type_, value)
	}
	ErrMoreThanLimit = func(type_ DeductionType, value float64) error {
		return fmt.Errorf("the %s deduction cannot be more than %f", type_, value)
	}
	ErrUpdateDatabase = func(type_ DeductionType) error {
		return fmt.Errorf("failed to set %s deduction", type_)
	}
)

func limiter(type_ DeductionType, amount, lower, upper float64) error {
	if amount < lower {
		return ErrLessThanLimit(type_, lower)
	}

	if amount > upper {
		return ErrMoreThanLimit(type_, upper)
	}

	return nil
}
