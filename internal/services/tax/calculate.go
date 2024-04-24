package tax

import (
	"context"
)

type CalculateRequest struct {
	Income     float64
	WHT        float64
	Allowances []Allowance
}

type CalculateResponse struct {
	Tax    float64
	Refund float64
}

func (s *service) Calculate(ctx context.Context, req CalculateRequest) (*CalculateResponse, error) {
	if req.Income < 0 {
		s.log.Fields(map[string]interface{}{"income": req.Income}).E("Income cannot be negative")
		return nil, ErrNegativeIncome
	}

	totalAllowances, err := s.calculateAllowances(req.Allowances)
	if err != nil {
		return nil, err
	}

	netIncome := max(req.Income-totalAllowances, 0)
	totalTax, err := calculateProgressiveTax(netIncome)
	if err != nil {
		s.log.Fields(map[string]interface{}{
			"netIncome":       netIncome,
			"income":          req.Income,
			"totalAllowances": totalAllowances,
		}).Err(err).E("Failed to calculate progressive tax")
		return nil, err
	}

	return &CalculateResponse{
		Tax:    max(totalTax-req.WHT, 0),
		Refund: max(req.WHT-totalTax, 0),
	}, nil
}
