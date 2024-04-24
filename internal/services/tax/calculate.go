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

const DefaultPersonalAllowances = 60000.0

func (s *service) Calculate(ctx context.Context, req CalculateRequest) (*CalculateResponse, error) {
	result, err := calculateProgressiveTax(req.Income - DefaultPersonalAllowances)
	if err != nil {
		return nil, err
	}

	return &CalculateResponse{
		Tax: result,
	}, nil
}
