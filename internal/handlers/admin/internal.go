package admin

import (
	"fmt"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request")
	ErrDeductPersonal = fmt.Errorf("unable to set personal deduction")
	ErrDeductKReceipt = fmt.Errorf("unable to set k-receipt deduction")
)

func toErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}
