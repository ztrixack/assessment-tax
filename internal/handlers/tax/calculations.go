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
		Allowances: remapAllowances(r.Allowances),
	}
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType" validate:"required,oneof=donation k-receipt" example:"donation"`
	Amount        float64 `json:"amount" validate:"min=0" example:"0.0"`
}

type CalculationsResponse struct {
	Tax       float64    `json:"tax"`
	TaxLevel  []TaxLevel `json:"taxLevel"`
	TaxRefund *float64   `json:"taxRefund,omitempty"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

// Calculations calculates the tax based on total income, withholding tax (WHT), and allowances.
//
//	@summary		Calculate Tax
//	@description	This endpoint calculates the tax and potentially applicable tax refund and tax levels based on the provided total income, withholding tax, and allowances.
//	@tags			tax
//	@accept			json
//	@produce		json
//	@param			request	body		CalculationsRequest		true	"Input request for tax calculation"
//	@success		200		{object}	CalculationsResponse	"Successfully calculated tax and returns the tax details"
//	@failure		400		{object}	ErrorResponse			"Bad request if the input validation fails"
//	@failure		500		{object}	ErrorResponse			"Internal server error if the tax calculations service fails"
//	@router			/tax/calculations [post]
func (h *handler) Calculations(c api.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if c.Request().Body == http.NoBody {
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	var req CalculationsRequest
	if err := c.Bind(&req); err != nil {
		h.log.Err(err).E("Failed to bind request")
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	if err := c.Validate(&req); err != nil {
		h.log.Err(err).Fields(logger.Fields{"request": req}).E("Failed to validate request")
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	res, err := h.tax.Calculate(ctx, req.toServiceRequest())
	if err != nil {
		h.log.Err(err).E("Failed to calculate tax")
		return c.JSON(http.StatusInternalServerError, toErrorResponse(ErrCalculateTax))
	}

	return c.JSON(http.StatusOK, toCalculationsResponse(*res))
}
