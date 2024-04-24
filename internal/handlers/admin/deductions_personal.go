package admin

import (
	"net/http"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
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

	res, err := dumpService(req)
	if err != nil {
		h.log.Err(err).E("Failed to set personal deduction")
		return c.JSON(http.StatusInternalServerError, toErrorResponse(ErrDeductPersonal))
	}

	return c.JSON(http.StatusOK, DeductionsPersonalResponse{PersonalDeduction: res})
}

func dumpService(req DeductionsPersonalRequest) (float64, error) {
	return req.Amount, nil
}
