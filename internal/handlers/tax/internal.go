package tax

import (
	"fmt"

	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrCalculateTax   = fmt.Errorf("failed to calculate tax")
)

func toCalculationsResponse(r tax.CalculateResponse) CalculationsResponse {
	return CalculationsResponse{
		Tax:       r.Tax,
		TaxRefund: remapTaxRefund(r.Refund),
	}
}

func toErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}

func remapTaxRefund(refund float64) *float64 {
	if refund == 0.0 {
		return nil
	}

	return &refund
}
