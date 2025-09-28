package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/thevibhu/plusgrade-tax-calculator/internal/models"
)

type TaxService interface {
	GetTaxBrackets(year string) ([]models.TaxBracket, error)
	CalculateTax(income float64, year string) (*models.TaxCalculationResponse, error)
}

type taxService struct {
	apiURL     string
	httpClient *http.Client
}

func NewTaxService(apiURL string) TaxService {
	return &taxService{
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type APIErrorDetail struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIErrorResponse is the structured error response from the interview test server API
type APIErrorResponse struct {
	Errors []APIErrorDetail `json:"errors"`
}

// Error implements the error interface for APIErrorResponse
func (e *APIErrorResponse) Error() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling API error response")
		return `{"errors":[{"code":"INTERNAL_ERROR","message":"failed to marshal api error response"}]}`
	}

	return string(jsonData)
}

func (s *taxService) GetTaxBrackets(year string) ([]models.TaxBracket, error) {
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%s", s.apiURL, year)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching tax brackets")
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Error().Err(readErr).Int("status", resp.StatusCode).Msg("API returned an error but failed to read response body")
			return nil, fmt.Errorf("API returned non-OK status: %d, and response body could not be read", resp.StatusCode)
		}

		var apiError APIErrorResponse
		if json.Unmarshal(body, &apiError) == nil && len(apiError.Errors) > 0 {
			log.Error().Err(&apiError).Msg("API returned structured errors")
			return nil, &apiError
		}

		log.Error().Int("status", resp.StatusCode).RawJSON("body", body).Msg("Received unparseable error response from API")
		return nil, fmt.Errorf("API error with status %d: %s", resp.StatusCode, string(body))
	}

	// Handle successful response
	var response models.TaxBracketsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error().Err(err).Msg("Error decoding successful response")
		return nil, err
	}

	log.Info().Msgf("Successfully fetched tax brackets for year %s", year)
	return response.TaxBrackets, nil
}

func (s *taxService) CalculateTax(income float64, year string) (*models.TaxCalculationResponse, error) {
	// Get tax brackets
	brackets, err := s.GetTaxBrackets(year)
	if err != nil {
		log.Error().Msgf("Error getting tax brackets: %v", err)
		return nil, err
	}

	if len(brackets) == 0 {
		return nil, fmt.Errorf("no tax brackets found for year %s", year)
	}

	// Calculate taxes
	totalTax := 0.0
	taxesByBand := []models.BandTaxDetail{}
	remainingIncome := income

	for _, bracket := range brackets {
		if remainingIncome <= 0 {
			break
		}

		bandDetail := models.BandTaxDetail{
			Min:  bracket.Min,
			Max:  bracket.Max,
			Rate: bracket.Rate,
		}

		// Calculate taxable income in this bracket
		var taxableInBracket float64
		if bracket.Max == 0 { // Last bracket (no upper limit)
			taxableInBracket = remainingIncome
		} else {
			bracketWidth := bracket.Max - bracket.Min
			if remainingIncome >= bracketWidth {
				taxableInBracket = bracketWidth
			} else {
				taxableInBracket = remainingIncome
			}
		}

		// Calculate tax for this bracket
		bandTax := taxableInBracket * bracket.Rate
		bandDetail.TaxableIncome = taxableInBracket
		bandDetail.TaxAmount = math.Round(bandTax*100) / 100

		taxesByBand = append(taxesByBand, bandDetail)
		totalTax += bandTax

		if bracket.Max > 0 {
			remainingIncome -= (bracket.Max - bracket.Min)
		} else {
			remainingIncome = 0
		}
	}

	// Calculate effective rate
	effectiveRate := 0.0
	if income > 0 {
		effectiveRate = (totalTax / income) * 100
	}

	response := &models.TaxCalculationResponse{
		Income:         income,
		TaxYear:        year,
		TotalTax:       math.Round(totalTax*100) / 100,
		TaxesByBand:    taxesByBand,
		EffectiveRate:  math.Round(effectiveRate*100) / 100,
		AfterTaxIncome: math.Round((income-totalTax)*100) / 100,
	}

	log.Info().Msgf("Tax calculation completed for income %.2f in year %s: total tax %.2f",
		income, year, response.TotalTax)

	return response, nil
}
