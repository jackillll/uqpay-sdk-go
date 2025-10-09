package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestBalances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		// Test getting balance for a specific currency
		currency := "USD"

		resp, err := client.Banking.Balances.Get(ctx, currency)
		if err != nil {
			t.Logf("❌ Get balance for %s returned error: %v", currency, err)
			return
		}

		t.Logf("✅ Balance retrieved for %s", currency)
		t.Logf("💰 Balance ID: %s", resp.BalanceID)
		t.Logf("   Currency: %s", resp.Currency)
		t.Logf("   Available: %s", resp.AvailableBalance)
		t.Logf("   Prepaid: %s", resp.PrepaidBalance)
		t.Logf("   Margin: %s", resp.MarginBalance)
		t.Logf("   Frozen: %s", resp.FrozenBalance)
		t.Logf("   Status: %s", resp.BalanceStatus)
		t.Logf("   Created: %s", resp.CreateTime)
		t.Logf("   Updated: %s", resp.UpdateTime)
	})

	t.Run("List", func(t *testing.T) {
		req := &banking.ListBalancesRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Balances.List(ctx, req)
		if err != nil {
			t.Logf("❌ List balances returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d balances (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			for i, balance := range resp.Data {
				t.Logf("💰 Balance %d: %s - Available: %s, Status: %s",
					i+1, balance.Currency, balance.AvailableBalance, balance.BalanceStatus)
			}
		} else {
			t.Logf("ℹ️  No balances found")
		}
	})

	t.Run("ListTransactions", func(t *testing.T) {
		req := &banking.ListBalanceTransactionsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Balances.ListTransactions(ctx, req)
		if err != nil {
			t.Logf("❌ List balance transactions returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d balance transactions (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			txn := resp.Data[0]
			t.Logf("🔍 First transaction:")
			t.Logf("   ID: %s", txn.TransactionID)
			t.Logf("   Type: %s", txn.TransactionType)
			t.Logf("   Amount: %s %s", txn.Amount, txn.Currency)
			t.Logf("   Status: %s", txn.TransactionStatus)
			t.Logf("   Balance Before: %s", txn.BalanceBefore)
			t.Logf("   Balance After: %s", txn.BalanceAfter)
			t.Logf("   Description: %s", txn.Description)
			t.Logf("   Reference ID: %s", txn.ReferenceID)
			t.Logf("   Created: %s", txn.CreateTime)
		} else {
			t.Logf("ℹ️  No balance transactions found")
		}
	})

	t.Run("ListTransactionsWithFilters", func(t *testing.T) {
		req := &banking.ListBalanceTransactionsRequest{
			PageSize:          10,
			PageNumber:        1,
			Currency:          "USD",
			TransactionType:   "DEPOSIT",
			TransactionStatus: "COMPLETED",
		}

		resp, err := client.Banking.Balances.ListTransactions(ctx, req)
		if err != nil {
			t.Logf("❌ List balance transactions with filters returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d completed USD deposit transactions (total: %d)",
			len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, txn := range resp.Data {
				t.Logf("💰 Transaction %d: %s %s, Type: %s, Balance: %s -> %s",
					i+1, txn.Amount, txn.Currency, txn.TransactionType,
					txn.BalanceBefore, txn.BalanceAfter)
			}
		}
	})

	t.Run("ListTransactionsByType", func(t *testing.T) {
		transactionTypes := []string{"PAYIN", "DEPOSIT", "PAYOUT", "TRANSFER", "CONVERSION", "FEE"}

		for _, txnType := range transactionTypes {
			req := &banking.ListBalanceTransactionsRequest{
				PageSize:        10,
				PageNumber:      1,
				TransactionType: txnType,
			}

			resp, err := client.Banking.Balances.ListTransactions(ctx, req)
			if err != nil {
				t.Logf("❌ List %s transactions returned error: %v", txnType, err)
				continue
			}

			t.Logf("📊 %s transactions: %d found", txnType, resp.TotalItems)
		}
	})
}
