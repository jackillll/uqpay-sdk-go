package banking

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	UpdateTime   string `json:"update_time"`
}

// ListRatesRequest represents a request to list exchange rates
type ListRatesRequest struct {
	CurrencyPairs []string `json:"currency_pairs,omitempty"` // optional: filter by specific currency pairs (e.g., ["USD/EUR", "GBP/USD"])
}

// ListRatesResponse represents a response containing exchange rates
type ListRatesResponse struct {
	Data struct {
		LastUpdated              time.Time  `json:"last_updated"`
		Rates                    []RateItem `json:"rates"`
		UnavailableCurrencyPairs []string   `json:"unavailable_currency_pairs"`
	} `json:"data"`
}

// List retrieves current exchange rates
// Optionally filter by specific currency pairs
func (c *ExchangeRatesClient) List(ctx context.Context, req *ListRatesRequest) (*ListRatesResponse, error) {
	var resp ListRatesResponse
	path := "/v1/exchange/rates"

	// Add currency_pairs query parameter if specified
	if req != nil && len(req.CurrencyPairs) > 0 {
		// Build dual-format list: original and no-slash
		seen := make(map[string]struct{})
		var out []string
		for _, p := range req.CurrencyPairs {
			if _, ok := seen[p]; !ok {
				seen[p] = struct{}{}
				out = append(out, p)
			}
			if strings.Contains(p, "/") {
				ns := strings.ReplaceAll(p, "/", "")
				if _, ok := seen[ns]; !ok {
					seen[ns] = struct{}{}
					out = append(out, ns)
				}
			}
		}
		pairs := strings.Join(out, ",")
		path += fmt.Sprintf("?currency_pairs=%s", pairs)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list exchange rates: %w", err)
	}

	return &resp, nil
}
