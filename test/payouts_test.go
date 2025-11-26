package test

import (
	"context"
	"testing"

	"github.com/jackillll/uqpay-sdk-go/banking"
)

func TestPayouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &banking.ListPayoutsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Payouts.List(ctx, req)
		if err != nil {
			t.Logf("❌ List payouts returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d payouts (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			payout := resp.Data[0]
			t.Logf("🔍 First payout:")
			t.Logf("   ID: %s", payout.PayoutID)
			t.Logf("   Reference: %s", payout.ShortReferenceID)
			t.Logf("   Amount: %s %s", payout.Amount, payout.Currency)
			t.Logf("   Fee: %s", payout.Fee)
			t.Logf("   Status: %s", payout.PayoutStatus)
			t.Logf("   Purpose: %s", payout.PayoutPurpose)
			t.Logf("   Beneficiary: %s", payout.Beneficiary.BeneficiaryName)
			t.Logf("   Created: %s", payout.CreateTime)
			if payout.CompletedTime != "" {
				t.Logf("   Completed: %s", payout.CompletedTime)
			}
			if payout.FailureReason != "" {
				t.Logf("   Failure Reason: %s", payout.FailureReason)
			}
		} else {
			t.Logf("ℹ️  No payouts found")
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		req := &banking.ListPayoutsRequest{
			PageSize:     10,
			PageNumber:   1,
			PayoutStatus: "COMPLETED",
			Currency:     "USD",
		}

		resp, err := client.Banking.Payouts.List(ctx, req)
		if err != nil {
			t.Logf("❌ List payouts with filters returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d completed USD payouts (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, payout := range resp.Data {
				t.Logf("💰 Payout %d: %s %s to %s, Status: %s",
					i+1, payout.Amount, payout.Currency,
					payout.Beneficiary.BeneficiaryName, payout.PayoutStatus)
			}
		}
	})

	t.Run("ListByStatus", func(t *testing.T) {
		statuses := []string{"PENDING", "PROCESSING", "COMPLETED", "FAILED", "CANCELLED"}

		for _, status := range statuses {
			req := &banking.ListPayoutsRequest{
				PageSize:     10,
				PageNumber:   1,
				PayoutStatus: status,
			}

			resp, err := client.Banking.Payouts.List(ctx, req)
			if err != nil {
				t.Logf("❌ List %s payouts returned error: %v", status, err)
				continue
			}

			t.Logf("📊 %s payouts: %d found", status, resp.TotalItems)
		}
	})

	t.Run("Create", func(t *testing.T) {
		// First, check if we have any beneficiaries
		beneficiariesReq := &banking.ListBeneficiariesRequest{
			PageSize:   10,
			PageNumber: 1,
			Status:     "active",
		}

		beneficiariesResp, err := client.Banking.Beneficiaries.List(ctx, beneficiariesReq)
		if err != nil {
			t.Logf("⏭️  Unable to list beneficiaries: %v", err)
			t.Skip("Skipping create test - cannot verify beneficiary availability")
		}

		if len(beneficiariesResp.Data) == 0 {
			t.Skip("⏭️  No active beneficiaries available - skipping create test")
		}

		// Skip actual creation to avoid test costs
		t.Skip("⏭️  Skipping actual payout creation to avoid transaction costs")

		// Uncomment to test actual payout creation
		/*
			beneficiaryID := beneficiariesResp.Data[0].BeneficiaryID
			t.Logf("🔍 Using beneficiary: %s", beneficiaryID)

			req := &banking.CreatePayoutRequest{
				BeneficiaryID: beneficiaryID,
				Currency:      "USD",
				Amount:        "10.00",
				PayoutPurpose: "test_payment",
				Description:   "Test payout via SDK",
				Reference:     "TEST-REF-001",
			}

			resp, err := client.Banking.Payouts.Create(ctx, req)
			if err != nil {
				t.Fatalf("❌ Failed to create payout: %v", err)
			}

			t.Logf("✅ Payout created successfully")
			t.Logf("💰 Payout ID: %s", resp.PayoutID)
			t.Logf("   Reference: %s", resp.ShortReferenceID)
			t.Logf("   Status: %s", resp.Status)
			t.Logf("   Created: %s", resp.CreateTime)
		*/
	})

	t.Run("CreateWithInlineBeneficiary", func(t *testing.T) {
		t.Skip("⏭️  Skipping inline beneficiary payout creation to avoid transaction costs")

		// Uncomment to test creating payout with inline beneficiary details
		/*
			req := &banking.CreatePayoutRequest{
				Currency:      "USD",
				Amount:        "10.00",
				PayoutPurpose: "vendor_payment",
				Description:   "Test payout with inline beneficiary",
				Beneficiary: &banking.PayoutBeneficiary{
					BeneficiaryName: "Test Beneficiary",
					BankDetails: &banking.PayoutBankDetails{
						AccountNumber: "1234567890",
						AccountName:   "Test Account",
						BankCode:      "TEST001",
						BankName:      "Test Bank",
						AccountType:   "CHECKING",
					},
					ContactDetails: &banking.PayoutContactDetails{
						Email:       "test@example.com",
						PhoneNumber: "+1234567890",
						Country:     "US",
					},
				},
			}

			resp, err := client.Banking.Payouts.Create(ctx, req)
			if err != nil {
				t.Fatalf("❌ Failed to create payout: %v", err)
			}

			t.Logf("✅ Payout created successfully")
			t.Logf("💰 Payout ID: %s", resp.PayoutID)
		*/
	})

	t.Run("Get", func(t *testing.T) {
		// First, list payouts to get a valid ID
		listReq := &banking.ListPayoutsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		listResp, err := client.Banking.Payouts.List(ctx, listReq)
		if err != nil {
			t.Logf("❌ Failed to list payouts: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Skip("⏭️  No payouts available to test Get operation")
		}

		payoutID := listResp.Data[0].PayoutID
		t.Logf("🔍 Getting payout details for ID: %s", payoutID)

		resp, err := client.Banking.Payouts.Get(ctx, payoutID)
		if err != nil {
			t.Fatalf("❌ Failed to get payout: %v", err)
		}

		t.Logf("✅ Payout retrieved successfully")
		t.Logf("💰 Payout ID: %s", resp.PayoutID)
		t.Logf("   Reference: %s", resp.ShortReferenceID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Fee: %s", resp.Fee)
		t.Logf("   Status: %s", resp.PayoutStatus)
		t.Logf("   Purpose: %s", resp.PayoutPurpose)
		t.Logf("   Beneficiary ID: %s", resp.Beneficiary.BeneficiaryID)
		t.Logf("   Beneficiary Name: %s", resp.Beneficiary.BeneficiaryName)

		if resp.Beneficiary.BankDetails != nil {
			t.Logf("   Bank: %s (%s)",
				resp.Beneficiary.BankDetails.BankName,
				resp.Beneficiary.BankDetails.AccountNumber)
		}

		if resp.Beneficiary.WalletDetails != nil {
			t.Logf("   Wallet: %s (%s)",
				resp.Beneficiary.WalletDetails.WalletProvider,
				resp.Beneficiary.WalletDetails.WalletNumber)
		}

		t.Logf("   Created: %s", resp.CreateTime)
		if resp.CompletedTime != "" {
			t.Logf("   Completed: %s", resp.CompletedTime)
		}
		if resp.FailureReason != "" {
			t.Logf("   Failure Reason: %s", resp.FailureReason)
		}

		if resp.TransactionDetails != nil {
			t.Logf("   Transaction ID: %s", resp.TransactionDetails.TransactionID)
			if resp.TransactionDetails.ExchangeRate != "" {
				t.Logf("   Exchange Rate: %s", resp.TransactionDetails.ExchangeRate)
			}
		}
	})
}
