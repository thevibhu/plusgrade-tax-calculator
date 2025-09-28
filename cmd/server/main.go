package main

import (
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thevibhu/plusgrade-tax-calculator/config"
	"github.com/thevibhu/plusgrade-tax-calculator/internal/handler"
	"github.com/thevibhu/plusgrade-tax-calculator/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Echo server
	e := echo.New()

	// Echo middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize service layer
	taxService := service.NewTaxService(cfg.TaxAPIURL)

	// Initialize handlers
	taxHandler := handler.NewTaxHandler(taxService)

	// Routes
	e.GET("/health", healthCheck)
	e.POST("/tax/calculate", taxHandler.CalculateTax)
	e.GET("/tax/brackets/:year", taxHandler.GetTaxBrackets)

	// Start server
	log.Info().Msgf("Starting server on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "healthy"})
}
