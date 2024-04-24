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
	if err := request.validate(); err != nil {
		s.log.Err(err).
			Fields(logger.Fields{"type": request.Type, "amount": request.Amount}).
			E("Invalid request to set %s deduction", request.Type)
		return 0, err
	}

	query := fmt.Sprintf("UPDATE allowances SET %s = $1", sanitizeType(request.Type))
	_, err := s.db.Execute(query, request.Amount)
	if err != nil {
		s.log.Err(err).
			Fields(logger.Fields{"query": query, "amount": request.Amount}).
			E("Failed to update %s deduction to allowances table in database", request.Type)
		return 0, ErrUpdateDatabase(request.Type)
	}

	return request.Amount, nil
}
