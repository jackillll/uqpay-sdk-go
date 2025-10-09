package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// TransactionsClient handles transaction operations
type TransactionsClient struct {
	client *common.APIClient
}

// Transaction represents a transaction
type Transaction struct {
	TransactionID       string `json:"transaction_id"`
	CardID              string `json:"card_id"`
	TransactionType     string `json:"transaction_type"`
	TransactionAmount   string `json:"transaction_amount"` // API returns string, not float
	TransactionCurrency string `json:"transaction_currency"`
	BillingAmount       string `json:"billing_amount"` // API returns string, not float
	BillingCurrency     string `json:"billing_currency"`
	MerchantName        string `json:"merchant_name"`
	TransactionStatus   string `json:"transaction_status"`
	TransactionTime     string `json:"transaction_time"`
}

// ListTransactionsRequest represents a transaction list request
type ListTransactionsRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	CardID     string `json:"card_id,omitempty"`
}

// ListTransactionsResponse represents a transaction list response
type ListTransactionsResponse struct {
	TotalPages int           `json:"total_pages"`
	TotalItems int           `json:"total_items"`
	Data       []Transaction `json:"data"`
}

// Get retrieves a transaction by ID
func (c *TransactionsClient) Get(ctx context.Context, transactionID string) (*Transaction, error) {
	var transaction Transaction
	path := fmt.Sprintf("/v1/issuing/transactions/%s", transactionID)
	if err := c.client.Get(ctx, path, &transaction); err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

// List lists transactions
func (c *TransactionsClient) List(ctx context.Context, req *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	var resp ListTransactionsResponse
	path := fmt.Sprintf("/v1/issuing/transactions?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if req.CardID != "" {
		path = fmt.Sprintf("%s&card_id=%s", path, req.CardID)
	}
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return &resp, nil
}
