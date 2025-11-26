package issuing

import (
	"context"
	"fmt"

	"github.com/jackillll/uqpay-sdk-go/common"
)

// CardholdersClient handles cardholder operations
type CardholdersClient struct {
	client *common.APIClient
}

// CreateCardholderRequest represents a cardholder creation request
type CreateCardholderRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CountryCode string `json:"country_code"`
}

// Cardholder represents a cardholder
type Cardholder struct {
	CardholderID string `json:"cardholder_id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	CountryCode  string `json:"country_code"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
}

// ListCardholdersRequest represents a cardholder list request
type ListCardholdersRequest struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}

// ListCardholdersResponse represents a cardholder list response
type ListCardholdersResponse struct {
	TotalPages int          `json:"total_pages"`
	TotalItems int          `json:"total_items"`
	Data       []Cardholder `json:"data"`
}

// Create creates a new cardholder
func (c *CardholdersClient) Create(ctx context.Context, req *CreateCardholderRequest) (*Cardholder, error) {
	var cardholder Cardholder
	if err := c.client.Post(ctx, "/v1/issuing/cardholders", req, &cardholder); err != nil {
		return nil, fmt.Errorf("failed to create cardholder: %w", err)
	}
	return &cardholder, nil
}

// Get retrieves a cardholder by ID
func (c *CardholdersClient) Get(ctx context.Context, cardholderID string) (*Cardholder, error) {
	var cardholder Cardholder
	path := fmt.Sprintf("/v1/issuing/cardholders/%s", cardholderID)
	if err := c.client.Get(ctx, path, &cardholder); err != nil {
		return nil, fmt.Errorf("failed to get cardholder: %w", err)
	}
	return &cardholder, nil
}

// List lists cardholders
func (c *CardholdersClient) List(ctx context.Context, req *ListCardholdersRequest) (*ListCardholdersResponse, error) {
	var resp ListCardholdersResponse
	path := fmt.Sprintf("/v1/issuing/cardholders?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list cardholders: %w", err)
	}
	return &resp, nil
}
