package tax

import (
	"context"
	"errors"
	"math"
)

type Servicer interface {
	Calculate(ctx context.Context, req CalculateRequest) (*CalculateResponse, error)
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

func calculateProgressiveTax(income float64) (float64, error) {
	if income < 0 {
		return 0, ErrNegativeIncome
	}

	total := 0.0
	brackets := []struct {
		lower float64
		upper float64
		rate  float64
	}{
		{0, 150000, 0},                   // 0% for income between 0 - 150,000
		{150000, 500000, 0.10},           // 10% for income between 150,001 - 500,000
		{500000, 1000000, 0.15},          // 15% for income between 500,001 - 1,000,000
		{1000000, 2000000, 0.20},         // 20% for income between 1,000,001 - 2,000,000
		{2000000, math.MaxFloat64, 0.35}, // 35% for income over 2,000,001
	}

	for _, bracket := range brackets {
		if income > bracket.lower {
			upper := math.Min(income, bracket.upper)
			total += calculateStepTax(income, bracket.lower, upper, bracket.rate)
		}
	}

	return total, nil
}
