package tax

import (
	"fmt"

	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrCalculateTax   = fmt.Errorf("failed to calculate tax")

	TaxLevelLabels = []string{"0-150,000", "150,001-500,000", "500,001-1,000,000", "1,000,001-2,000,000", "2,000,001 ขึ้นไป"}
)

func toCalculationsResponse(r tax.CalculateResponse) CalculationsResponse {
	return CalculationsResponse{
		Tax:       r.Tax,
		TaxLevel:  remapTaxLevel(TaxLevelLabels, r.TaxLevel),
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

func remapAllowances(allowances []Allowance) []tax.Allowance {
	result := make([]tax.Allowance, len(allowances))

	for i, a := range allowances {
		result[i] = tax.Allowance{
			Type:   tax.AllowanceType(a.AllowanceType),
			Amount: a.Amount,
		}
	}

	return result
}

func remapTaxLevel(labels []string, levels []float64) []TaxLevel {
	result := make([]TaxLevel, len(levels))

	if len(levels) != len(labels) {
		labels = make([]string, len(levels))
		for i := range labels {
			labels[i] = fmt.Sprintf("Bucket: #%d", i+1)
		}
	}

	for i := range levels {
		result[i] = TaxLevel{
			Level: labels[i],
			Tax:   levels[i],
		}
	}

	return result
}
