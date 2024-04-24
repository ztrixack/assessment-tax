package tax

import (
	"fmt"

	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrCalculateTax   = fmt.Errorf("failed to calculate tax")
)

func toResponse(r tax.CalculateResponse) CalculationsResponse {
	return CalculationsResponse{
		Tax: r.Tax,
	}
}
