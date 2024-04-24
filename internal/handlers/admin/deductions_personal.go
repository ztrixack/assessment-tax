package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/admin"
)

type DeductionsPersonalRequest struct {
	Amount float64 `json:"amount" validate:"min=10000,max=100000" example:"60000.0"`
}

type DeductionsPersonalResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

// DeductionsPersonal sets a personal deduction by admin.
//
//	@summary		Set personal deduction
//	@description	Sets the personal deduction based on the provided request parameters.
//	@tags			admin/deductions
//	@accept			json
//	@produce		json
//	@param			request	body	DeductionsPersonalRequest	true	"Input request for setting personal deduction"
//	@security		BasicAuth
//	@success		200	{object}	DeductionsPersonalResponse	"Successfully response with updated deduction details"
//	@failure		400	{object}	ErrorResponse				"Bad request if the input validation fails"
//	@failure		401	{object}	ErrorResponse				"Unauthorized"
//	@failure		500	{object}	ErrorResponse				"Internal Server Error if there is a problem setting the deduction"
//	@router			/admin/deductions/personal [post]
func (h handler) DeductionsPersonal(c api.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if c.Request().Body == http.NoBody {
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	var req DeductionsPersonalRequest
	if err := c.Bind(&req); err != nil {
		h.log.Err(err).E("Failed to bind request")
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	if err := c.Validate(&req); err != nil {
		h.log.Err(err).Fields(logger.Fields{"request": req}).E("Failed to validate request")
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	res, err := h.admin.SetDeduction(ctx, req.toServiceRequest())
	if err != nil {
		h.log.Err(err).E("Failed to set personal deduction")
		return c.JSON(http.StatusInternalServerError, toErrorResponse(ErrDeductPersonal))
	}

	return c.JSON(http.StatusOK, DeductionsPersonalResponse{PersonalDeduction: res})
}

func (r *DeductionsPersonalRequest) toServiceRequest() admin.SetDeductionRequest {
	return admin.SetDeductionRequest{
		Type:   admin.Personal,
		Amount: r.Amount,
	}
}
