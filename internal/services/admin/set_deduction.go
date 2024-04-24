package admin

import (
	"context"
	"fmt"

	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

type SetDeductionRequest struct {
	Type   DeductionType `json:"type"`
	Amount float64       `json:"amount"`
}

func (s *service) SetDeduction(ctx context.Context, request SetDeductionRequest) (float64, error) {
	var typeStr string
	if request.Type == Personal {
		typeStr = "personal"
		if err := limiter(Personal, request.Amount, PersonalMinimum, PersonalMaximum); err != nil {
			s.log.Err(err).E("Invalid %s deduction amount", request.Type)
			return 0, err
		}
	} else if request.Type == KReceipt {
		typeStr = "k_receipt"
		if err := limiter(KReceipt, request.Amount, KReceiptMinimum, KReceiptMaximum); err != nil {
			s.log.Err(err).E("Invalid %s deduction amount", request.Type)
			return 0, err
		}
	} else {
		return 0, ErrInvalidDeductionType
	}

	query := fmt.Sprintf("UPDATE allowances SET %s = $1", typeStr)
	_, err := s.db.Execute(query, request.Amount)
	if err != nil {
		s.log.Err(err).
			Fields(logger.Fields{"query": query, "amount": request.Amount}).
			E("Failed to update %s deduction to allowances table in database", request.Type)
		return 0, ErrUpdateDatabase(request.Type)
	}

	return request.Amount, nil
}
