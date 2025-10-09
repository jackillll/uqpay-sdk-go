package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// ConversionClient handles conversion operations
type ConversionClient struct {
	client *common.APIClient
}

// Conversion represents a currency conversion
type Conversion struct {
	ConversionID     string `json:"conversion_id"`
	ShortReferenceID string `json:"short_reference_id"`
	CurrencyFrom     string `json:"currency_from"`
	CurrencyTo       string `json:"currency_to"`
	AmountFrom       string `json:"amount_from"`
	AmountTo         string `json:"amount_to"`
	Rate             string `json:"rate"`
	ConversionStatus string `json:"conversion_status"` // COMPLETED, PENDING, FAILED
	CreateTime       string `json:"create_time"`
	CompletedTime    string `json:"completed_time,omitempty"`
	SettlementDate   string `json:"settlement_date,omitempty"`
}

// CreateConversionRequest represents a conversion creation request
type CreateConversionRequest struct {
	CurrencyFrom   string `json:"currency_from"`   // required
	CurrencyTo     string `json:"currency_to"`     // required
	AmountFrom     string `json:"amount_from"`     // required
	SettlementDate string `json:"settlement_date"` // optional, format: YYYY-MM-DD
	QuoteID        string `json:"quote_id"`        // optional, if provided, conversion will use quoted rate
}

// CreateConversionResponse represents a conversion creation response
type CreateConversionResponse struct {
	ConversionID     string `json:"conversion_id"`
	ShortReferenceID string `json:"short_reference_id"`
}

// ListConversionsRequest represents a conversion list request
type ListConversionsRequest struct {
	PageSize         int    `json:"page_size"`         // required, 10-100
	PageNumber       int    `json:"page_number"`       // required, >=1
	StartTime        string `json:"start_time"`        // optional, ISO8601
	EndTime          string `json:"end_time"`          // optional, ISO8601
	ConversionStatus string `json:"conversion_status"` // optional: COMPLETED, PENDING, FAILED
	CurrencyFrom     string `json:"currency_from"`     // optional
	CurrencyTo       string `json:"currency_to"`       // optional
}

// ListConversionsResponse represents a conversion list response
type ListConversionsResponse struct {
	TotalPages int          `json:"total_pages"`
	TotalItems int          `json:"total_items"`
	Data       []Conversion `json:"data"`
}

// CreateQuoteRequest represents a quote creation request
type CreateQuoteRequest struct {
	CurrencyFrom   string `json:"currency_from"`   // required
	CurrencyTo     string `json:"currency_to"`     // required
	AmountFrom     string `json:"amount_from"`     // required
	SettlementDate string `json:"settlement_date"` // optional, format: YYYY-MM-DD
}

// CreateQuoteResponse represents a quote creation response
type CreateQuoteResponse struct {
	QuoteID        string `json:"quote_id"`
	CurrencyFrom   string `json:"currency_from"`
	CurrencyTo     string `json:"currency_to"`
	AmountFrom     string `json:"amount_from"`
	AmountTo       string `json:"amount_to"`
	Rate           string `json:"rate"`
	SettlementDate string `json:"settlement_date,omitempty"`
	ExpiresAt      string `json:"expires_at"` // ISO8601 timestamp
}

// ConversionDate represents available conversion dates for a currency pair
type ConversionDate struct {
	Date          string `json:"date"`           // format: YYYY-MM-DD
	FirstCutoff   string `json:"first_cutoff"`   // ISO8601 timestamp
	SecondCutoff  string `json:"second_cutoff"`  // ISO8601 timestamp
	OptimizedDate bool   `json:"optimized_date"` // whether this is the optimal conversion date
}

// List lists conversions
func (c *ConversionClient) List(ctx context.Context, req *ListConversionsRequest) (*ListConversionsResponse, error) {
	var resp ListConversionsResponse
	path := fmt.Sprintf("/v1/conversion?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.StartTime != "" {
		path += fmt.Sprintf("&start_time=%s", req.StartTime)
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("&end_time=%s", req.EndTime)
	}
	if req.ConversionStatus != "" {
		path += fmt.Sprintf("&conversion_status=%s", req.ConversionStatus)
	}
	if req.CurrencyFrom != "" {
		path += fmt.Sprintf("&currency_from=%s", req.CurrencyFrom)
	}
	if req.CurrencyTo != "" {
		path += fmt.Sprintf("&currency_to=%s", req.CurrencyTo)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list conversions: %w", err)
	}
	return &resp, nil
}

// Create creates a new conversion
func (c *ConversionClient) Create(ctx context.Context, req *CreateConversionRequest) (*CreateConversionResponse, error) {
	var resp CreateConversionResponse
	if err := c.client.Post(ctx, "/v1/conversion", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create conversion: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific conversion
func (c *ConversionClient) Get(ctx context.Context, conversionID string) (*Conversion, error) {
	var resp Conversion
	path := fmt.Sprintf("/v1/conversion/%s", conversionID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get conversion: %w", err)
	}
	return &resp, nil
}

// ListConversionDates retrieves available conversion dates for a currency pair
func (c *ConversionClient) ListConversionDates(ctx context.Context, currencyFrom, currencyTo string) ([]ConversionDate, error) {
	var resp []ConversionDate
	path := fmt.Sprintf("/v1/conversion/conversion_dates?currency_from=%s&currency_to=%s", currencyFrom, currencyTo)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list conversion dates: %w", err)
	}
	return resp, nil
}

// CreateQuote creates a new conversion quote
func (c *ConversionClient) CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*CreateQuoteResponse, error) {
	var resp CreateQuoteResponse
	if err := c.client.Post(ctx, "/v1/conversion/quote", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}
	return &resp, nil
}
