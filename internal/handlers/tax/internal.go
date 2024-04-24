package tax

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrCalculateTax   = fmt.Errorf("failed to calculate tax")
	ErrInvalidFile    = fmt.Errorf("invalid file")

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

func getFileFromRequest(c api.Context) (*multipart.FileHeader, error) {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func parseCSVFile(_ *multipart.FileHeader) ([]tax.CalculateRequest, error) {
	// TODO
	return []tax.CalculateRequest{
		{
			Income:     500000,
			WHT:        0,
			Allowances: []tax.Allowance{{Type: "donation", Amount: 0}},
		},
		{
			Income:     600000,
			WHT:        40000,
			Allowances: []tax.Allowance{{Type: "donation", Amount: 20000}},
		},
		{
			Income:     750000,
			WHT:        50000,
			Allowances: []tax.Allowance{{Type: "donation", Amount: 15000}},
		},
	}, nil
}

func (h *handler) calculateTaxes(ctx context.Context, reqs []tax.CalculateRequest) ([]Tax, error) {
	taxes := make([]Tax, 0, len(reqs))
	for _, req := range reqs {
		res, err := h.tax.Calculate(ctx, req)
		if err != nil {
			return nil, err
		}
		taxes = append(taxes, toTax(req.Income, *res))
	}
	return taxes, nil
}

func toTax(income float64, r tax.CalculateResponse) Tax {
	result := Tax{
		TotalIncome: income,
		Tax:         r.Tax,
	}

	if r.Refund > 0 {
		result.TaxRefund = &r.Refund
	}

	return result
}
