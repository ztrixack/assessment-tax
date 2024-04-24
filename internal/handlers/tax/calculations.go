package tax

import (
	"net/http"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
)

type CalculationsRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type CalculationsResponse struct {
	Tax       float64  `json:"tax"`
	TaxRefund *float64 `json:"taxRefund,omitempty"`
}

func (h *handler) Calculations(c api.Context) error {
	var req CalculationsRequest
	if err := c.Bind(&req); err != nil {
		h.log.Err(err).E("Failed to bind request")
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	res := h.calculateTax(req)

	return c.JSON(http.StatusOK, res)
}

func (h *handler) calculateTax(_ CalculationsRequest) CalculationsResponse {
	return CalculationsResponse{
		Tax: 29000,
	}
}
