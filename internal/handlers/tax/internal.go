package tax

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
	"github.com/ztrixack/assessment-tax/internal/utils/csv"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrCalculateTax   = fmt.Errorf("failed to calculate tax")
	ErrInvalidFile    = fmt.Errorf("invalid file")
	ErrGetFileFailed  = fmt.Errorf("failed to get CSV file")

	TaxLevelLabels = []string{"0-150,000", "150,001-500,000", "500,001-1,000,000", "1,000,001-2,000,000", "2,000,001 ขึ้นไป"}
)

type ErrorResponse struct {
	Error string `json:"error"`
}

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

func parseCSVFile(file *multipart.FileHeader) ([]tax.CalculateRequest, error) {
	csvReader, fileCloser, err := csv.OpenCSV(file)
	if err != nil {
		return nil, err
	}
	defer fileCloser.Close()

	expectedHeaders := []string{"totalIncome", "wht", "donation"}
	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	if err := csv.ValidateHeaders(headers, expectedHeaders); err != nil {
		return nil, err
	}

	interfaceResults, err := csv.ProcessRecords(csvReader, parseTaxRecord)
	if err != nil {
		return nil, err
	}

	var results []tax.CalculateRequest
	for _, ir := range interfaceResults {
		if tr, ok := ir.(*tax.CalculateRequest); ok {
			results = append(results, *tr)
		} else {
			return nil, fmt.Errorf("type assertion failed")
		}
	}

	return results, nil
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

func parseTaxRecord(record []string) (interface{}, error) {
	if len(record) != 3 {
		return nil, io.ErrUnexpectedEOF
	}

	totalIncome, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return nil, err
	}

	wht, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}

	donation, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, err
	}

	return &tax.CalculateRequest{
		Income:     totalIncome,
		WHT:        wht,
		Allowances: []tax.Allowance{{Type: "donation", Amount: donation}},
	}, nil
}
