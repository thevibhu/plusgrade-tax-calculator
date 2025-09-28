package models

// TaxBracket represents a single tax bracket
type TaxBracket struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max,omitempty"`
	Rate float64 `json:"rate"`
}

// TaxBracketsResponse from external API
type TaxBracketsResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}

// TaxCalculationRequest request
type TaxCalculationRequest struct {
	Income  float64 `json:"income" validate:"required"`
	TaxYear string  `json:"tax_year" validate:"required"`
}

// TaxCalculationResponse response
type TaxCalculationResponse struct {
	Income         float64         `json:"income"`
	TaxYear        string          `json:"tax_year"`
	TotalTax       float64         `json:"total_tax"`
	TaxesByBand    []BandTaxDetail `json:"taxes_by_band"`
	EffectiveRate  float64         `json:"effective_rate"`
	AfterTaxIncome float64         `json:"after_tax_income"`
}

// BandTaxDetail for each tax bracket
type BandTaxDetail struct {
	Min           float64 `json:"min"`
	Max           float64 `json:"max,omitempty"`
	Rate          float64 `json:"rate"`
	TaxableIncome float64 `json:"taxable_income"`
	TaxAmount     float64 `json:"tax_amount"`
}
