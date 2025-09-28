# Plusgrade Tax Calculator API

**Vaibhav's Plusgrade Technical Assessment**

An income tax calculation service built with Go, featuring robust error handling, comprehensive testing, and production ready architecture.

## Setup & Run

```bash

# Clone the repo
git clone https://github.com/vaibhavverma/plusgrade-tax-calculator.git
cd plusgrade-tax-calculator

# Build and run the service (you should have docker installed)
docker compose up -d --build

# Clean environment after use
docker compose down --volumes --remove-orphans

# API will be available at http://localhost:8080
# Interview test server runs on http://localhost:5001
```

## Architecture & Design

This implementation follows **Clean Architecture** principles:

```
├── cmd/server/          # Server entry point
├── config/              # Config management
├── internal/
│   ├── handler/         # Rest handler
│   ├── service/         # Service layer
│   └── models/          # Domain structs
├── docker/              # Deployment and testing Dockerfiles
├── mocks/               # Generated unit test mocks
└── postman/             # Integration test suite
```

### Key Design Patterns

- **Dependency Injection**: Services are injected into handlers, enabling easy testing and modularity
- **Repository Pattern**: Clean abstraction over external API calls
- **Error Wrapping**: Structured error handling with proper HTTP status codes


## Core Features

### Tax Calculation
- **Accurate marginal tax rate calculations** for Canadian federal tax brackets
- **Multi-bracket support** with proper income distribution across bands
- **Precision handling** with proper rounding to 2 decimal places
- **Edge case coverage** including zero income scenarios

### Robust Error Handling
- **Graceful degradation** when external API fails randomly (as per assessment requirements)
- **Structured error responses** that mirror external API format
- **Comprehensive input validation** with clear error messages
- **Timeout protection** with 10-second HTTP client timeouts

### Production-Ready Features
- **Health check endpoint** for load balancer integration
- **Structured logging** with zerolog for observability
- **CORS support** for cross-origin requests
- **Graceful shutdown** handling

## API Endpoints

### Calculate Tax
```http
POST /tax/calculate
Content-Type: application/json

{
  "income": 100000,
  "tax_year": "2022"
}
```

**Response:**
```json
{
  "income": 100000,
  "tax_year": "2022",
  "total_tax": 17739.17,
  "taxes_by_band": [
    {
      "min": 0,
      "max": 50197,
      "rate": 0.15,
      "taxable_income": 50197,
      "tax_amount": 7529.55
    },
    {
      "min": 50197,
      "max": 100392,
      "rate": 0.205,
      "taxable_income": 49803,
      "tax_amount": 10209.62
    }
  ],
  "effective_rate": 17.74,
  "after_tax_income": 82260.83
}
```

### Get Tax Brackets
```http
GET /tax/brackets/2022
```
**Response:**
```json
{
    "tax_brackets": [
        {
            "min": 0,
            "max": 50197,
            "rate": 0.15
        },
        {
            "min": 50197,
            "max": 100392,
            "rate": 0.205
        },
        {
            "min": 100392,
            "max": 155625,
            "rate": 0.26
        },
        {
            "min": 155625,
            "max": 221708,
            "rate": 0.29
        },
        {
            "min": 221708,
            "rate": 0.33
        }
    ]
}
```


### Health Check
```http
GET /health
```

**Response:**
```json
{
    "status": "healthy"
}
```

## Testing Strategy

### Multi-Layered Testing Approach

**Unit Tests** (`internal/service/tax_service_test.go`)
- Tax calculation logic validation
- Error handling scenarios
- Mock external API responses
- Edge case coverage (zero income, high income)

**Integration Tests** (Postman Collection)
- End-to-end API validation
- Real service interaction
- Error response verification

### Running Tests

```bash
# If needed to regenerate mocks (you need to have mockery and Go installed)
# Install Go if not already installed: https://go.dev/doc/install
# Install mockery
go install github.com/vektra/mockery/v2@latest 

# Regenerate mocks
mockery --all 

# Unit tests with coverage
docker-compose --profile test-suite up unit-test

# Integration tests with newman
docker-compose --profile test-suite up integration-test
```

### Test Coverage Highlights
- **Zero income calculation**: Validates $0 tax on $0 income
- **Single bracket income**: $50,000 → $7,500 tax (15% rate)
- **Multi-bracket income**: $100,000 → $17,739.17 tax (17.74% effective rate)
- **High income validation**: $1,234,567 → $385,587.65 tax (31.24% effective rate)
- **Error handling**: Invalid years, negative income, API failures
- **API integration**: Structured error responses, timeout handling

## Tech Stack

**Core Framework**
- **Go 1.19+** - Performance and reliability
- **Echo v4** - Lightweight, high-performance HTTP framework
- **Zerolog** - Structured, high-performance logging

**Testing & Quality**
- **Testify** - Comprehensive testing utilities
- **Mockery** - Auto-generated mocks for interfaces
- **Postman/Newman** - API integration testing

**Infrastructure**
- **Docker & Docker Compose** - Containerized deployment
- **Alpine Linux** - Minimal, secure base images
- **Multi-stage builds** - Optimized container sizes

## Error Handling

### External API Integration
The service gracefully handles the interview test server's intentional failures:

```go
// Structured error response matching external API format
type APIErrorResponse struct {
    Errors []APIErrorDetail `json:"errors"`
}

// Proper error wrapping and classification
if errors.As(err, &apiError) {
    return c.JSON(http.StatusBadGateway, apiError)
}
```

### Input Validation
- **Income validation**: Non-negative values only
- **Year validation**: Supports 2019-2022 as per API constraints
- **Request format validation**: Proper JSON structure enforcement

## Logging, Monitoring & Scalability

### Structured Logging
```go
log.Info().Msgf("Tax calculation completed for income %.2f in year %s: total tax %.2f",
    income, year, response.TotalTax)
```

### Health Monitoring
- `/health` endpoint returns service status
- Docker health checks with configurable intervals
- Proper HTTP status codes for monitoring systems


### Scalability Features
- **Stateless design**
- **Configuration via environment variables**
- **Health checks**

### Code Quality
- **gofmt** and **go vet** compliant
- **Dependency injection** for testability
- **Interface-driven design** for modularity
- **Clean Architecture** principles

**Status**: All requirements implemented
---