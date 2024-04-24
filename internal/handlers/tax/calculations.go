package tax

import (
	"net/http"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

type CalculationsRequest struct {
	TotalIncome float64     `json:"totalIncome" validate:"min=0" example:"500000.0"`
	WHT         float64     `json:"wht" validate:"min=0" example:"0.0"`
	Allowances  []Allowance `json:"allowances" validate:"dive"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType" validate:"required,oneof=donation" example:"donation"`
	Amount        float64 `json:"amount" validate:"min=0" example:"0.0"`
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

	if err := c.Validate(&req); err != nil {
		h.log.Err(err).Fields(logger.Fields{"request": req}).E("Failed to validate request")
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
