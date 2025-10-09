package connect

import "github.com/uqpay/uqpay-sdk-go/common"

// Client provides access to Connect APIs
type Client struct {
	Accounts *AccountsClient
}

// NewClient creates a new Connect client
func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		Accounts: &AccountsClient{client: apiClient},
	}
}
