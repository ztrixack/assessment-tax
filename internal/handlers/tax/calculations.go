package tax

import (
	"context"
	"net/http"
	"time"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/tax"
)

type CalculationsRequest struct {
	TotalIncome float64     `json:"totalIncome" validate:"min=0" example:"500000.0"`
	WHT         float64     `json:"wht" validate:"min=0" example:"0.0"`
	Allowances  []Allowance `json:"allowances" validate:"dive"`
}

func (r *CalculationsRequest) toServiceRequest() tax.CalculateRequest {
	return tax.CalculateRequest{
		Income:     r.TotalIncome,
		WHT:        r.WHT,
		Allowances: []tax.Allowance{},
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var req CalculationsRequest
	if err := c.Bind(&req); err != nil {
		h.log.Err(err).E("Failed to bind request")
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	if err := c.Validate(&req); err != nil {
		h.log.Err(err).Fields(logger.Fields{"request": req}).E("Failed to validate request")
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	res := h.tax.Calculate(ctx, req.toServiceRequest())

	return c.JSON(http.StatusOK, toResponse(*res))
}
