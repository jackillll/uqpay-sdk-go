package banking

import (
	"context"
	"fmt"
	"strings"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// ExchangeRatesClient handles exchange rate operations
type ExchangeRatesClient struct {
	client *common.APIClient
}

// RateItem represents an exchange rate for a currency pair
type RateItem struct {
	CurrencyPair string `json:"currency_pair"` // e.g., "USD/EUR"
	BuyPrice     string `json:"buy_price"`     // Price for buying the base currency
	SellPrice    string `json:"sell_price"`    // Price for selling the base currency
	UpdateTime   string `json:"update_time"`   // Timestamp of the rate
}

// ListRatesRequest represents a request to list exchange rates
type ListRatesRequest struct {
	CurrencyPairs []string `json:"currency_pairs,omitempty"` // optional: filter by specific currency pairs (e.g., ["USD/EUR", "GBP/USD"])
}

// ListRatesResponse represents a response containing exchange rates
type ListRatesResponse struct {
	Rates []RateItem `json:"rates"`
}

// List retrieves current exchange rates
// Optionally filter by specific currency pairs
func (c *ExchangeRatesClient) List(ctx context.Context, req *ListRatesRequest) (*ListRatesResponse, error) {
	var resp ListRatesResponse
	path := "/v1/exchange/rates"

	// Add currency_pairs query parameter if specified
	if req != nil && len(req.CurrencyPairs) > 0 {
		pairs := strings.Join(req.CurrencyPairs, ",")
		path += fmt.Sprintf("?currency_pairs=%s", pairs)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list exchange rates: %w", err)
	}
	return &resp, nil
}
