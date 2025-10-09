package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestDeposits(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &banking.ListDepositsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Deposits.List(ctx, req)
		if err != nil {
			t.Logf("❌ List deposits returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d deposits (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			deposit := resp.Data[0]
			t.Logf("🔍 First deposit:")
			t.Logf("   ID: %s", deposit.DepositID)
			t.Logf("   Reference: %s", deposit.ShortReferenceID)
			t.Logf("   Amount: %s %s", deposit.Amount, deposit.Currency)
			t.Logf("   Status: %s", deposit.DepositStatus)
			t.Logf("   Payment Method: %s", deposit.PaymentMethod)
			t.Logf("   Payer: %s (%s)", deposit.PayerName, deposit.PayerEmail)
			t.Logf("   Description: %s", deposit.Description)
			t.Logf("   Created: %s", deposit.CreateTime)
			if deposit.CompletedTime != "" {
				t.Logf("   Completed: %s", deposit.CompletedTime)
			}
		} else {
			t.Logf("ℹ️  No deposits found")
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		req := &banking.ListDepositsRequest{
			PageSize:      10,
			PageNumber:    1,
			DepositStatus: "COMPLETED",
			Currency:      "USD",
		}

		resp, err := client.Banking.Deposits.List(ctx, req)
		if err != nil {
			t.Logf("❌ List deposits with filters returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d completed USD deposits (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, deposit := range resp.Data {
				t.Logf("💰 Deposit %d: %s %s from %s, Status: %s",
					i+1, deposit.Amount, deposit.Currency,
					deposit.PayerName, deposit.DepositStatus)
			}
		}
	})

	t.Run("ListByStatus", func(t *testing.T) {
		statuses := []string{"PENDING", "COMPLETED", "FAILED"}

		for _, status := range statuses {
			req := &banking.ListDepositsRequest{
				PageSize:      10,
				PageNumber:    1,
				DepositStatus: status,
			}

			resp, err := client.Banking.Deposits.List(ctx, req)
			if err != nil {
				t.Logf("❌ List %s deposits returned error: %v", status, err)
				continue
			}

			t.Logf("📊 %s deposits: %d found", status, resp.TotalItems)
		}
	})

	t.Run("ListByCurrency", func(t *testing.T) {
		currencies := []string{"USD", "EUR", "GBP", "KES", "UGX"}

		for _, currency := range currencies {
			req := &banking.ListDepositsRequest{
				PageSize:   10,
				PageNumber: 1,
				Currency:   currency,
			}

			resp, err := client.Banking.Deposits.List(ctx, req)
			if err != nil {
				t.Logf("❌ List %s deposits returned error: %v", currency, err)
				continue
			}

			t.Logf("📊 %s deposits: %d found", currency, resp.TotalItems)
		}
	})

	t.Run("ListWithTimeRange", func(t *testing.T) {
		// Example: List deposits from the last 30 days
		// Note: Adjust these dates as needed for your testing
		req := &banking.ListDepositsRequest{
			PageSize:   10,
			PageNumber: 1,
			// StartTime: "2024-01-01T00:00:00Z",  // Uncomment and adjust as needed
			// EndTime:   "2024-01-31T23:59:59Z",  // Uncomment and adjust as needed
		}

		resp, err := client.Banking.Deposits.List(ctx, req)
		if err != nil {
			t.Logf("❌ List deposits with time range returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d deposits in time range (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("Get", func(t *testing.T) {
		// First, list deposits to get a valid ID
		listReq := &banking.ListDepositsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		listResp, err := client.Banking.Deposits.List(ctx, listReq)
		if err != nil {
			t.Logf("❌ Failed to list deposits: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Skip("⏭️  No deposits available to test Get operation")
		}

		depositID := listResp.Data[0].DepositID
		t.Logf("🔍 Getting deposit details for ID: %s", depositID)

		resp, err := client.Banking.Deposits.Get(ctx, depositID)
		if err != nil {
			t.Fatalf("❌ Failed to get deposit: %v", err)
		}

		t.Logf("✅ Deposit retrieved successfully")
		t.Logf("💰 Deposit ID: %s", resp.DepositID)
		t.Logf("   Reference: %s", resp.ShortReferenceID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.DepositStatus)
		t.Logf("   Payment Method: %s", resp.PaymentMethod)
		t.Logf("   Payer Name: %s", resp.PayerName)
		t.Logf("   Payer Email: %s", resp.PayerEmail)
		t.Logf("   Description: %s", resp.Description)
		t.Logf("   Created: %s", resp.CreateTime)
		if resp.CompletedTime != "" {
			t.Logf("   Completed: %s", resp.CompletedTime)
		}
	})

	t.Run("GetMultipleDeposits", func(t *testing.T) {
		// List deposits
		listReq := &banking.ListDepositsRequest{
			PageSize:   5,
			PageNumber: 1,
		}

		listResp, err := client.Banking.Deposits.List(ctx, listReq)
		if err != nil {
			t.Logf("❌ Failed to list deposits: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Skip("⏭️  No deposits available")
		}

		// Get details for each deposit
		t.Logf("🔍 Retrieving details for %d deposits", len(listResp.Data))

		for i, deposit := range listResp.Data {
			resp, err := client.Banking.Deposits.Get(ctx, deposit.DepositID)
			if err != nil {
				t.Logf("❌ Failed to get deposit %s: %v", deposit.DepositID, err)
				continue
			}

			t.Logf("💰 Deposit %d: %s %s - %s (Status: %s)",
				i+1, resp.Amount, resp.Currency,
				resp.PayerName, resp.DepositStatus)
		}
	})
}
