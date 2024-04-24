package tax

import "context"

type CalculateRequest struct {
	Income     float64
	WHT        float64
	Allowances []Allowance
}

type CalculateResponse struct {
	Tax float64
}

func (s *service) Calculate(ctx context.Context, req CalculateRequest) (*CalculateResponse, error) {
	if req.Income < 0 {
		return nil, ErrNegativeIncome
	}

	totalAllowances, err := s.calculateAllowances(req.Allowances)
	if err != nil {
		return nil, err
	}

	netIncome := max(req.Income-totalAllowances, 0)
	result, err := calculateProgressiveTax(netIncome)
	if err != nil {
		return nil, err
	}

	return &CalculateResponse{Tax: result}, nil
}
