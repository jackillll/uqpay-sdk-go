package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestVirtualAccountsCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	// Create virtual account
	req := &banking.CreateVirtualAccountRequest{
		VirtualAccountName: "Test Virtual Account",
		Currencies:         []string{"USD", "EUR"},
	}

	account, err := client.Banking.VirtualAccounts.Create(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}

	if account.VirtualAccountID == "" {
		t.Error("Expected virtual_account_id to be set")
	}
	if account.VirtualAccountName != req.VirtualAccountName {
		t.Errorf("Expected virtual_account_name %s, got %s", req.VirtualAccountName, account.VirtualAccountName)
	}
	if account.Status == "" {
		t.Error("Expected status to be set")
	}
	if account.CreateTime == "" {
		t.Error("Expected create_time to be set")
	}
	if len(account.CurrencyBankDetail) == 0 {
		t.Error("Expected currency_bank_detail to have entries")
	}

	t.Logf("Created virtual account: %s", account.VirtualAccountID)
	t.Logf("  Name: %s", account.VirtualAccountName)
	t.Logf("  Status: %s", account.Status)
	t.Logf("  Create Time: %s", account.CreateTime)
	t.Logf("  Currency Bank Details: %d", len(account.CurrencyBankDetail))

	for i, detail := range account.CurrencyBankDetail {
		t.Logf("  Detail %d:", i+1)
		t.Logf("    Currency: %s", detail.Currency)
		t.Logf("    Bank Name: %s", detail.BankName)
		t.Logf("    Bank Address: %s", detail.BankAddress)
		t.Logf("    Bank Country Code: %s", detail.BankCountryCode)
		t.Logf("    Account Name: %s", detail.AccountName)
		t.Logf("    Account Number: %s", detail.AccountNumber)
		if detail.SwiftCode != "" {
			t.Logf("    Swift Code: %s", detail.SwiftCode)
		}
		if detail.RoutingNumber != "" {
			t.Logf("    Routing Number: %s", detail.RoutingNumber)
		}
		if detail.IBAN != "" {
			t.Logf("    IBAN: %s", detail.IBAN)
		}
		if detail.SortCode != "" {
			t.Logf("    Sort Code: %s", detail.SortCode)
		}
		if detail.IFSCCode != "" {
			t.Logf("    IFSC Code: %s", detail.IFSCCode)
		}
	}
}

func TestVirtualAccountsList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	req := &banking.ListVirtualAccountsRequest{
		PageSize:   10,
		PageNumber: 1,
	}

	resp, err := client.Banking.VirtualAccounts.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list virtual accounts: %v", err)
	}

	if resp.TotalPages < 0 {
		t.Error("Expected total_pages to be >= 0")
	}
	if resp.TotalItems < 0 {
		t.Error("Expected total_items to be >= 0")
	}

	t.Logf("Listed virtual accounts:")
	t.Logf("  Total Pages: %d", resp.TotalPages)
	t.Logf("  Total Items: %d", resp.TotalItems)
	t.Logf("  Current Page Items: %d", len(resp.Data))

	for i, account := range resp.Data {
		t.Logf("  Virtual Account %d:", i+1)
		t.Logf("    ID: %s", account.VirtualAccountID)
		t.Logf("    Name: %s", account.VirtualAccountName)
		t.Logf("    Status: %s", account.Status)
		t.Logf("    Create Time: %s", account.CreateTime)
		t.Logf("    Currency Bank Details: %d", len(account.CurrencyBankDetail))
	}
}

func TestVirtualAccountsCreateAndList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	// Create a virtual account
	createReq := &banking.CreateVirtualAccountRequest{
		VirtualAccountName: "Integration Test Account",
		Currencies:         []string{"USD"},
	}

	account, err := client.Banking.VirtualAccounts.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}

	t.Logf("Created virtual account: %s", account.VirtualAccountID)

	// List virtual accounts and verify the created account is in the list
	listReq := &banking.ListVirtualAccountsRequest{
		PageSize:   10,
		PageNumber: 1,
	}

	listResp, err := client.Banking.VirtualAccounts.List(ctx, listReq)
	if err != nil {
		t.Fatalf("Failed to list virtual accounts: %v", err)
	}

	found := false
	for _, acc := range listResp.Data {
		if acc.VirtualAccountID == account.VirtualAccountID {
			found = true
			if acc.VirtualAccountName != createReq.VirtualAccountName {
				t.Errorf("Expected account name %s, got %s", createReq.VirtualAccountName, acc.VirtualAccountName)
			}
			break
		}
	}

	if !found {
		t.Error("Created virtual account not found in list")
	} else {
		t.Log("Successfully verified created account in list")
	}
}

func TestVirtualAccountsMultipleCurrencies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	// Create virtual account with multiple currencies
	req := &banking.CreateVirtualAccountRequest{
		VirtualAccountName: "Multi-Currency Test Account",
		Currencies:         []string{"USD", "EUR", "GBP"},
	}

	account, err := client.Banking.VirtualAccounts.Create(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}

	t.Logf("Created multi-currency virtual account: %s", account.VirtualAccountID)

	if len(account.CurrencyBankDetail) != len(req.Currencies) {
		t.Errorf("Expected %d currency bank details, got %d", len(req.Currencies), len(account.CurrencyBankDetail))
	}

	currencyMap := make(map[string]bool)
	for _, detail := range account.CurrencyBankDetail {
		currencyMap[detail.Currency] = true
		t.Logf("  Currency: %s, Bank: %s", detail.Currency, detail.BankName)
	}

	for _, currency := range req.Currencies {
		if !currencyMap[currency] {
			t.Errorf("Expected currency %s not found in bank details", currency)
		}
	}
}
