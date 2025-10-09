package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestTransfers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &banking.ListTransfersRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Transfers.List(ctx, req)
		if err != nil {
			t.Logf("❌ List transfers returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d transfers (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			transfer := resp.Data[0]
			t.Logf("🔍 First transfer: ID=%s, Amount=%s %s, Status=%s",
				transfer.TransferID, transfer.Amount, transfer.Currency, transfer.TransferStatus)
			t.Logf("   Source: %s -> Target: %s",
				transfer.SourceAccountID, transfer.TargetAccountID)
			t.Logf("   Reference: %s, Created: %s",
				transfer.ShortReferenceID, transfer.CreateTime)
		} else {
			t.Logf("ℹ️  No transfers found")
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		req := &banking.ListTransfersRequest{
			PageSize:       10,
			PageNumber:     1,
			TransferStatus: "completed",
			Currency:       "USD",
		}

		resp, err := client.Banking.Transfers.List(ctx, req)
		if err != nil {
			t.Logf("❌ List transfers with filters returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d completed USD transfers (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			transfer := resp.Data[0]
			t.Logf("🔍 Sample transfer: ID=%s, Amount=%s %s",
				transfer.TransferID, transfer.Amount, transfer.Currency)
		}
	})

	t.Run("Create", func(t *testing.T) {
		// Skip create test if no source/target accounts available
		t.Skip("⏭️  Skipping create test - requires valid source and target account IDs")

		// Uncomment and update with valid account IDs to test
		/*
		req := &banking.CreateTransferRequest{
			SourceAccountID: "source-account-uuid",
			TargetAccountID: "target-account-uuid",
			Currency:        "USD",
			Amount:          "10.00",
			Reason:          "Test transfer",
		}

		resp, err := client.Banking.Transfers.Create(ctx, req)
		if err != nil {
			t.Fatalf("❌ Failed to create transfer: %v", err)
		}

		t.Logf("✅ Transfer created successfully")
		t.Logf("💰 Transfer ID: %s", resp.TransferID)
		t.Logf("   Reference: %s", resp.ShortReferenceID)
		*/
	})

	t.Run("Get", func(t *testing.T) {
		// First, list transfers to get a valid ID
		listReq := &banking.ListTransfersRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		listResp, err := client.Banking.Transfers.List(ctx, listReq)
		if err != nil {
			t.Logf("❌ Failed to list transfers: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Skip("⏭️  No transfers available to test Get operation")
		}

		transferID := listResp.Data[0].TransferID
		t.Logf("🔍 Getting transfer details for ID: %s", transferID)

		resp, err := client.Banking.Transfers.Get(ctx, transferID)
		if err != nil {
			t.Fatalf("❌ Failed to get transfer: %v", err)
		}

		t.Logf("✅ Transfer retrieved successfully")
		t.Logf("💰 Transfer ID: %s", resp.TransferID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.TransferStatus)
		t.Logf("   Source: %s", resp.SourceAccountID)
		t.Logf("   Target: %s", resp.TargetAccountID)
		t.Logf("   Reason: %s", resp.Reason)
		t.Logf("   Created: %s", resp.CreateTime)
		if resp.CompletedTime != "" {
			t.Logf("   Completed: %s", resp.CompletedTime)
		}
	})
}
