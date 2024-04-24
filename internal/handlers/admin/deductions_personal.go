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
	Amount float64 `json:"amount" validate:"min=10000,max=100000" example:"0.0"`
}

type DeductionsPersonalResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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
