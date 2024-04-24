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

// UploadCSV handles the uploading and processing of a CSV file
//
//	@summary		Upload CSV file
//	@description	Uploads a CSV file and parses it to JSON.
//	@tags			tax
//	@accept			multipart/form-data
//	@produce		json
//	@param			taxFile	formData	file				true	"Upload CSV tax file"
//	@success		200		{object}	UploadCSVResponse	"Successfully parsed tax data"
//	@failure		400		{object}	ErrorResponse		"Unable to process the file, error in file retrieval or content"
//	@failure		500		{object}	ErrorResponse		"Internal server error, failed to read CSV header or records"
//	@router			/tax/calculations/upload-csv [post]
func (h *handler) UploadCSV(c api.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file, err := getFileFromRequest(c)
	if err != nil {
		h.log.Err(err).E("Failed to get file from request")
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrGetFileFailed))
	}

	reqs, err := parseCSVFile(file)
	if err != nil {
		h.log.Err(err).E("Failed to parse CSV file")
		return c.JSON(http.StatusBadRequest, toErrorResponse(ErrInvalidFile))
	}

	taxes, err := h.calculateTaxes(ctx, reqs)
	if err != nil {
		h.log.Err(err).E("Failed to calculate taxes")
		return c.JSON(http.StatusInternalServerError, toErrorResponse(ErrCalculateTax))
	}

	return c.JSON(http.StatusOK, UploadCSVResponse{Taxes: taxes})
}
