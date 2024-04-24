package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
	"github.com/ztrixack/assessment-tax/internal/services/admin"
)

type DeductionsKReceiptRequest struct {
	Amount float64 `json:"amount" validate:"min=0,max=100000" example:"50000.0"`
}

type DeductionsKReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

// DeductionsKReceipt sets a k-receipt deduction by admin.
//
//	@summary		Set k-receipt deduction
//	@description	Sets the k-receipt deduction based on the provided request parameters.
//	@tags			admin/deductions
//	@accept			json
//	@produce		json
//	@param			request	body	DeductionsKReceiptRequest	true	"Input request for setting k-receipt deduction"
//	@security		BasicAuth
//	@success		200			{object}		DeductionsKReceiptResponse	"Successfully response with updated deduction details"
//	@failure		{object}	ErrorResponse	400							"Bad request if the input validation fails"
//	@failure		{object}	ErrorResponse	401							"Unauthorized"
//	@failure		{object}	ErrorResponse	500							"Internal Server Error if there is a problem setting the deduction"
//	@router			/admin/deductions/k-receipt [post]
func (h handler) DeductionsKReceipt(c api.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if c.Request().Body == http.NoBody {
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidRequest))
	}

	var req DeductionsKReceiptRequest
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
		h.log.Err(err).E("Failed to set KReceipt deduction")
		return c.JSON(http.StatusInternalServerError, toErrorResponse(ErrDeductKReceipt))
	}

	return c.JSON(http.StatusOK, DeductionsKReceiptResponse{KReceipt: res})
}

func (r *DeductionsKReceiptRequest) toServiceRequest() admin.SetDeductionRequest {
	return admin.SetDeductionRequest{
		Type:   admin.KReceipt,
		Amount: r.Amount,
	}
}
