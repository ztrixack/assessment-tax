package tax

import (
	"context"
	"errors"
)

type Servicer interface {
	Calculate(ctx context.Context, req CalculateRequest) *CalculateResponse
}

type Allowance struct {
	Type   string
	Amount float64
}

var (
	ErrNegativeIncome = errors.New("income cannot be negative")
)

func calculateStepTax(income, lower, upper float64, rate float64) float64 {
	taxableIncome := min(income, upper) - lower
	return taxableIncome * rate
}

func calculateProgressiveTax(income float64) float64 {
	// 5 tax steps
	taxes := []float64{0, 0, 0, 0, 0}

	// range from 0 - 150,000, 0% tax
	taxes[0] = 0

	// range from 150,001 - 500,000, 10% tax
	if income > 150000 {
		taxes[1] = calculateStepTax(income, 150000, 500000, 0.1)
	}

	// range from 500,001 - 1,000,000, 15% tax
	if income > 500000 {
		taxes[2] = calculateStepTax(income, 500000, 1000000, 0.15)
	}

	// range from 1,000,001 - 2,000,000, 20% tax
	if income > 1000000 {
		taxes[3] = calculateStepTax(income, 1000000, 2000000, 0.2)
	}

	// range from 2,000,001 and above, 35% tax
	if income > 2000000 {
		taxes[4] = calculateStepTax(income, 2000000, income, 0.35)
	}

	sum := 0.0
	for _, tax := range taxes {
		sum += tax
	}

	return sum
}
