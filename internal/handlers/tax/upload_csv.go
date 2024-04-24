package tax

import (
	"context"
	"net/http"
	"time"

	"github.com/ztrixack/assessment-tax/internal/modules/api"
)

type UploadCSVResponse struct {
	Taxes []Tax `json:"taxes"`
}

type Tax struct {
	TotalIncome float64  `json:"totalIncome"`
	Tax         float64  `json:"tax"`
	TaxRefund   *float64 `json:"taxRefund,omitempty"`
}

func (h *handler) UploadCSV(c api.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file, err := getFileFromRequest(c)
	if err != nil {
		h.log.Err(err).E("Failed to get file from request")
		return c.JSON(http.StatusBadRequest, ErrInvalidFile)
	}

	reqs, err := parseCSVFile(file)
	if err != nil {
		h.log.Err(err).E("Failed to parse CSV file")
		return c.JSON(http.StatusBadRequest, ErrInvalidFile)
	}

	taxes, err := h.calculateTaxes(ctx, reqs)
	if err != nil {
		h.log.Err(err).E("Failed to calculate taxes")
		return c.JSON(http.StatusInternalServerError, ErrCalculateTax)
	}

	return c.JSON(http.StatusOK, UploadCSVResponse{Taxes: taxes})
}
