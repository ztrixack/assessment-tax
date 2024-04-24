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
	totalAllowances, err := s.calculateAllowances(req.Allowances)
	if err != nil {
		return nil, err
	}

	result, err := calculateProgressiveTax(req.Income - totalAllowances)
	if err != nil {
		return nil, err
	}

	return &CalculateResponse{Tax: result}, nil
}
