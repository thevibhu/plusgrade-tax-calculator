package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thevibhu/plusgrade-tax-calculator/internal/models"
)

// TestCalculateTax validates the tax calculation logic for various income levels
func TestCalculateTax(t *testing.T) {
	// Create a mock server to simulate the interview test server API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tax-calculator/tax-year/2022", r.URL.Path)
		response := models.TaxBracketsResponse{
			TaxBrackets: []models.TaxBracket{
				{Min: 0, Max: 50197, Rate: 0.15},
				{Min: 50197, Max: 100392, Rate: 0.205},
				{Min: 100392, Max: 155625, Rate: 0.26},
				{Min: 155625, Max: 221708, Rate: 0.29},
				{Min: 221708, Rate: 0.33},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service := NewTaxService(server.URL)

	// Test cases for successful calculations
	testCases := []struct {
		name         string
		income       float64
		expectedTax  float64
		expectedRate float64
	}{
		{"Zero income", 0, 0, 0},
		{"Income in first bracket", 50000, 7500.00, 15.00},
		{"Income spanning multiple brackets", 100000, 17739.17, 17.74},
		{"High income spanning all brackets", 1234567, 385587.65, 31.24},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.CalculateTax(tc.income, "2022")
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.InDelta(t, tc.expectedTax, result.TotalTax, 0.01)
			assert.InDelta(t, tc.expectedRate, result.EffectiveRate, 0.01)
		})
	}

	// Test case for when the interview test server API call fails
	t.Run("Error fetching tax brackets", func(t *testing.T) {
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errorServer.Close()

		errorService := NewTaxService(errorServer.URL)
		result, err := errorService.CalculateTax(50000, "2023")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "API error with status 500")
	})
}

// TestGetTaxBrackets validates the behavior of the API call to fetch tax brackets.
func TestGetTaxBrackets(t *testing.T) {
	t.Run("Success fetching brackets", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"tax_brackets":[{"min":0,"max":50197,"rate":0.15}]}`)
		}))
		defer server.Close()

		service := NewTaxService(server.URL)
		brackets, err := service.GetTaxBrackets("2022")

		assert.NoError(t, err)
		assert.NotNil(t, brackets)
		assert.Len(t, brackets, 1)
		assert.Equal(t, 0.15, brackets[0].Rate)
	})

	t.Run("API returns a generic server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal Server Error")
		}))
		defer server.Close()

		service := NewTaxService(server.URL)
		brackets, err := service.GetTaxBrackets("2022")

		assert.Error(t, err)
		assert.Nil(t, brackets)
		assert.Contains(t, err.Error(), "API error with status 500: Internal Server Error")
	})

	t.Run("API returns a structured error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			errorResp := APIErrorResponse{
				Errors: []APIErrorDetail{
					{Code: "INVALID_YEAR", Field: "year", Message: "Tax year not found"},
				},
			}
			json.NewEncoder(w).Encode(errorResp)
		}))
		defer server.Close()

		service := NewTaxService(server.URL)
		brackets, err := service.GetTaxBrackets("invalid-year")

		assert.Error(t, err)
		assert.Nil(t, brackets)

		// Assert that the error is the specific APIErrorResponse type
		apiErr, ok := err.(*APIErrorResponse)
		assert.True(t, ok, "error should be of type *APIErrorResponse")
		assert.Len(t, apiErr.Errors, 1)
		assert.Equal(t, "INVALID_YEAR", apiErr.Errors[0].Code)
	})

	t.Run("API returns OK status with malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"tax_brackets":[`) // Invalid JSON
		}))
		defer server.Close()

		service := NewTaxService(server.URL)
		brackets, err := service.GetTaxBrackets("2022")

		assert.Error(t, err)
		assert.Nil(t, brackets)
	})
}
