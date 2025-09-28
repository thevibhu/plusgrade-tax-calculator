package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thevibhu/plusgrade-tax-calculator/internal/models"
	"github.com/thevibhu/plusgrade-tax-calculator/internal/service"
)

type TaxHandler struct {
	taxService service.TaxService
}

func NewTaxHandler(taxService service.TaxService) *TaxHandler {
	return &TaxHandler{
		taxService: taxService,
	}
}

func (h *TaxHandler) CalculateTax(c echo.Context) error {
	var req models.TaxCalculationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	if req.Income < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Income must be non-negative",
		})
	}

	validYears := map[string]bool{"2019": true, "2020": true, "2021": true, "2022": true}
	if !validYears[req.TaxYear] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tax year must be one of: 2019, 2020, 2021, 2022",
		})
	}

	result, err := h.taxService.CalculateTax(req.Income, req.TaxYear)
	if err != nil {
		var apiError *service.APIErrorResponse

		// Check if the error is the specific API error type.
		if errors.As(err, &apiError) {
			return c.JSON(http.StatusBadGateway, apiError)
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "An internal server error occurred.",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaxHandler) GetTaxBrackets(c echo.Context) error {
	year := c.Param("year")

	validYears := map[string]bool{"2019": true, "2020": true, "2021": true, "2022": true}
	if !validYears[year] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tax year must be one of: 2019, 2020, 2021, 2022",
		})
	}

	brackets, err := h.taxService.GetTaxBrackets(year)
	if err != nil {
		var apiError *service.APIErrorResponse

		// Check if the error is the specific API error type.
		if errors.As(err, &apiError) {
			return c.JSON(http.StatusBadGateway, apiError)
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "An internal server error occurred.",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tax_brackets": brackets,
	})
}
