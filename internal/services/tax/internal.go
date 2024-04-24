package tax

import (
	"context"
	"fmt"
	"math"
)

type Servicer interface {
	Calculate(ctx context.Context, req CalculateRequest) (*CalculateResponse, error)
}

type Allowance struct {
	Type   AllowanceType
	Amount float64
}

type AllowanceList map[AllowanceType]float64

type AllowanceType string

const (
	Personal AllowanceType = "personal"
	Donation AllowanceType = "donation"
)

var (
	ErrNegativeIncome           = fmt.Errorf("income cannot be negative")
	ErrNegativeAllowanceAmount  = fmt.Errorf("allowance amount cannot be negative")
	ErrUnsupportedAllowanceType = fmt.Errorf("allowance type not supported")
)

func (s *service) calculateAllowances(allowanceList []Allowance) (float64, error) {
	allowances, err := s.getAllowances()
	if err != nil {
		s.log.Err(err).E("Failed to get allowances from database.")
		return 0, err
	}

	personal := allowances[Personal]
	donation := 0.0

	for _, allowance := range allowanceList {
		if allowance.Amount < 0 {
			s.log.Fields(map[string]interface{}{"allowance": allowance}).W("Allowance amount cannot be negative.")
			return 0, ErrNegativeAllowanceAmount
		}

		switch allowance.Type {
		case Donation:
			donation += allowance.Amount

		default:
			s.log.Fields(map[string]interface{}{"allowance": allowance}).W("Allowance type not supported.")
			return 0, ErrUnsupportedAllowanceType
		}
	}

	donation = calculateAllowance(donation, 0, allowances[Donation])

	return personal + donation, nil
}

func calculateAllowance(amount, lower, upper float64) float64 {
	return min(max(amount, lower), upper)
}

func (s *service) getAllowances() (AllowanceList, error) {
	row, err := s.db.QueryOne("SELECT personal, donation FROM allowances")
	if err != nil {
		return nil, err
	}

	var personal, donation float64
	err = row.Scan(&personal, &donation)
	if err != nil {
		return nil, err
	}

	return AllowanceList{Personal: personal, Donation: donation}, nil
}

func calculateStepTax(income, lower, upper float64, rate float64) float64 {
	taxableIncome := min(income, upper) - lower
	return taxableIncome * rate
}

func calculateProgressiveTax(income float64) (float64, []float64, error) {
	if income < 0 {
		return 0, nil, ErrNegativeIncome
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
	steps := make([]float64, len(brackets))

	for i, bracket := range brackets {
		if income > bracket.lower {
			upper := math.Min(income, bracket.upper)
			tax := calculateStepTax(income, bracket.lower, upper, bracket.rate)
			steps[i] = tax
			total += tax
		}
	}

	return total, steps, nil
}
