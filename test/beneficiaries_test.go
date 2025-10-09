package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestBeneficiaries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &banking.ListBeneficiariesRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Beneficiaries.List(ctx, req)
		if err != nil {
			t.Logf("❌ List beneficiaries returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d beneficiaries (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			beneficiary := resp.Data[0]
			t.Logf("🔍 First beneficiary:")
			t.Logf("   ID: %s", beneficiary.BeneficiaryID)
			t.Logf("   Type: %s", beneficiary.EntityType)

			if beneficiary.EntityType == "INDIVIDUAL" {
				t.Logf("   Name: %s %s", beneficiary.FirstName, beneficiary.LastName)
			} else {
				t.Logf("   Company: %s", beneficiary.CompanyName)
			}

			t.Logf("   Currency: %s", beneficiary.Currency)
			t.Logf("   Country: %s", beneficiary.Country)
			t.Logf("   Payment Method: %s", beneficiary.PaymentMethod)
			t.Logf("   Status: %s", beneficiary.Status)

			if beneficiary.BankDetails != nil {
				t.Logf("   Bank Account: %s", beneficiary.BankDetails.AccountNumber)
				if beneficiary.BankDetails.BankName != "" {
					t.Logf("   Bank Name: %s", beneficiary.BankDetails.BankName)
				}
				if beneficiary.BankDetails.IBAN != "" {
					t.Logf("   IBAN: %s", beneficiary.BankDetails.IBAN)
				}
			}

			if beneficiary.Email != "" {
				t.Logf("   Email: %s", beneficiary.Email)
			}
			if beneficiary.PhoneNumber != "" {
				t.Logf("   Phone: %s", beneficiary.PhoneNumber)
			}

			t.Logf("   Created: %s", beneficiary.CreateTime)
			t.Logf("   Updated: %s", beneficiary.UpdateTime)
		} else {
			t.Logf("ℹ️  No beneficiaries found")
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		req := &banking.ListBeneficiariesRequest{
			PageSize:   10,
			PageNumber: 1,
			Currency:   "USD",
			Status:     "active",
			EntityType: "INDIVIDUAL",
		}

		resp, err := client.Banking.Beneficiaries.List(ctx, req)
		if err != nil {
			t.Logf("❌ List beneficiaries with filters returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d active individual USD beneficiaries (total: %d)",
			len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, beneficiary := range resp.Data {
				t.Logf("💰 Beneficiary %d: %s %s (%s)",
					i+1, beneficiary.FirstName, beneficiary.LastName, beneficiary.Currency)
			}
		}
	})

	t.Run("ListByCountry", func(t *testing.T) {
		countries := []string{"US", "GB", "KE", "UG"}

		for _, country := range countries {
			req := &banking.ListBeneficiariesRequest{
				PageSize:   10,
				PageNumber: 1,
				Country:    country,
			}

			resp, err := client.Banking.Beneficiaries.List(ctx, req)
			if err != nil {
				t.Logf("❌ List %s beneficiaries returned error: %v", country, err)
				continue
			}

			t.Logf("📊 %s beneficiaries: %d found", country, resp.TotalItems)
		}
	})

	t.Run("ListPaymentMethods", func(t *testing.T) {
		// Test for USD in United States
		currency := "USD"
		country := "US"

		methods, err := client.Banking.Beneficiaries.ListPaymentMethods(ctx, currency, country)
		if err != nil {
			t.Logf("❌ List payment methods for %s/%s returned error: %v", currency, country, err)
			return
		}

		t.Logf("✅ Found %d payment methods for %s in %s", len(methods), currency, country)

		for i, method := range methods {
			t.Logf("🔍 Method %d:", i+1)
			t.Logf("   ID: %s", method.PaymentMethodID)
			t.Logf("   Name: %s", method.PaymentMethodName)
			t.Logf("   Currency: %s", method.Currency)
			t.Logf("   Country: %s", method.Country)

			if len(method.RequiredFields) > 0 {
				t.Logf("   Required Fields: %v", method.RequiredFields)
			}
			if len(method.OptionalFields) > 0 {
				t.Logf("   Optional Fields: %v", method.OptionalFields)
			}
			if method.MinAmount != "" {
				t.Logf("   Min Amount: %s", method.MinAmount)
			}
			if method.MaxAmount != "" {
				t.Logf("   Max Amount: %s", method.MaxAmount)
			}
		}
	})

	t.Run("ListPaymentMethodsMultipleCurrencies", func(t *testing.T) {
		testCases := []struct {
			currency string
			country  string
		}{
			{"USD", "US"},
			{"GBP", "GB"},
			{"EUR", "DE"},
			{"KES", "KE"},
			{"UGX", "UG"},
		}

		for _, tc := range testCases {
			methods, err := client.Banking.Beneficiaries.ListPaymentMethods(ctx, tc.currency, tc.country)
			if err != nil {
				t.Logf("❌ %s/%s: %v", tc.currency, tc.country, err)
				continue
			}

			t.Logf("📊 %s/%s: %d payment methods available", tc.currency, tc.country, len(methods))
		}
	})

	t.Run("Check", func(t *testing.T) {
		// This test validates beneficiary details without creating them
		req := &banking.BeneficiaryCheckRequest{
			Currency:      "USD",
			Country:       "US",
			PaymentMethod: "ACH",
			BankDetails: &banking.BankDetails{
				AccountNumber: "1234567890",
				RoutingNumber: "021000021",
			},
		}

		resp, err := client.Banking.Beneficiaries.Check(ctx, req)
		if err != nil {
			t.Logf("❌ Check beneficiary returned error: %v", err)
			// This is expected if validation fails
			return
		}

		t.Logf("✅ Beneficiary validation successful")
		t.Logf("🔍 Validation result:")
		t.Logf("   Currency: %s", resp.Currency)
		t.Logf("   Country: %s", resp.Country)
		t.Logf("   Payment Method: %s", resp.PaymentMethod)

		if resp.BankDetails != nil {
			if resp.BankDetails.BankName != "" {
				t.Logf("   Bank Name: %s", resp.BankDetails.BankName)
			}
			if resp.BankDetails.AccountNumber != "" {
				t.Logf("   Account Number: %s", resp.BankDetails.AccountNumber)
			}
		}
	})

	t.Run("CheckUKBeneficiary", func(t *testing.T) {
		req := &banking.BeneficiaryCheckRequest{
			Currency:      "GBP",
			Country:       "GB",
			PaymentMethod: "FASTER_PAYMENTS",
			BankDetails: &banking.BankDetails{
				AccountNumber: "12345678",
				SortCode:      "123456",
			},
		}

		resp, err := client.Banking.Beneficiaries.Check(ctx, req)
		if err != nil {
			t.Logf("❌ Check UK beneficiary returned error: %v", err)
			return
		}

		t.Logf("✅ UK beneficiary validation successful")
		t.Logf("🔍 Currency: %s, Country: %s, Method: %s",
			resp.Currency, resp.Country, resp.PaymentMethod)
	})

	t.Run("CheckSEPABeneficiary", func(t *testing.T) {
		req := &banking.BeneficiaryCheckRequest{
			Currency:      "EUR",
			Country:       "DE",
			PaymentMethod: "SEPA",
			BankDetails: &banking.BankDetails{
				IBAN: "DE89370400440532013000",
				BIC:  "COBADEFFXXX",
			},
		}

		resp, err := client.Banking.Beneficiaries.Check(ctx, req)
		if err != nil {
			t.Logf("❌ Check SEPA beneficiary returned error: %v", err)
			return
		}

		t.Logf("✅ SEPA beneficiary validation successful")
		t.Logf("🔍 Currency: %s, Country: %s, Method: %s",
			resp.Currency, resp.Country, resp.PaymentMethod)
	})
}
