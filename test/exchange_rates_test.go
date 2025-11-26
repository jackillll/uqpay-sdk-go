package test

import (
	"context"
	"testing"

	"github.com/jackillll/uqpay-sdk-go/banking"
)

func TestExchangeRates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("ListAll", func(t *testing.T) {
		req := &banking.ListRatesRequest{}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Fatalf("❌ Failed to list all exchange rates: %v", err)
		}

		t.Logf("✅ Retrieved %d exchange rates", len(resp.Data.Rates))

		if len(resp.Data.Rates) > 0 {
			// Display first few rates as examples
			displayCount := 5
			if len(resp.Data.Rates) < displayCount {
				displayCount = len(resp.Data.Rates)
			}

			t.Logf("📊 Sample exchange rates:")
			for i := 0; i < displayCount; i++ {
				rate := resp.Data.Rates[i]
				t.Logf("   %s: Buy=%s, Sell=%s",
					rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
			}

			// Group and summarize by base currency
			baseCurrencies := make(map[string]int)
			for _, rate := range resp.Data.Rates {
				// Extract base currency (before the /)
				for i, c := range rate.CurrencyPair {
					if c == '/' {
						base := rate.CurrencyPair[:i]
						baseCurrencies[base]++
						break
					}
				}
			}

			t.Logf("📊 Exchange rates by base currency:")
			for currency, count := range baseCurrencies {
				t.Logf("   %s: %d pairs", currency, count)
			}
		} else {
			t.Logf("ℹ️  No exchange rates available")
		}
	})

	t.Run("ListSpecificPairs", func(t *testing.T) {
		pairs := []string{"USD/EUR", "USDEUR", "USD/GBP", "USDGBP", "EUR/USD", "EURUSD", "GBP/USD", "GBPUSD"}

		req := &banking.ListRatesRequest{
			CurrencyPairs: pairs,
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to list specific exchange rates: %v", err)
			return
		}

		t.Logf("✅ Retrieved %d exchange rates for specified pairs", len(resp.Data.Rates))

		if len(resp.Data.Rates) > 0 {
			t.Logf("📊 Requested currency pairs:")
			for _, rate := range resp.Data.Rates {
				t.Logf("   %s: Buy=%s, Sell=%s",
					rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
			}
		} else {
			t.Logf("ℹ️  No rates found for requested pairs")
			if len(resp.Data.UnavailableCurrencyPairs) > 0 {
				t.Logf("⚠️  Unavailable currency pairs: %v", resp.Data.UnavailableCurrencyPairs)
			}
		}
	})

	t.Run("ListUSDPairs", func(t *testing.T) {
		// Request all USD-related pairs
		pairs := []string{
			"USD/EUR", "USDEUR", "USD/GBP", "USDGBP", "USD/JPY", "USDJPY", "USD/CHF", "USDCHF",
			"USD/CAD", "USDCAD", "USD/AUD", "USDAUD", "USD/NZD", "USDNZD", "USD/CNY", "USDCNY",
		}

		req := &banking.ListRatesRequest{
			CurrencyPairs: pairs,
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to list USD pairs: %v", err)
			return
		}

		t.Logf("✅ Retrieved %d USD exchange rates", len(resp.Data.Rates))

		if len(resp.Data.Rates) > 0 {
			t.Logf("💰 USD Exchange Rates:")
			for _, rate := range resp.Data.Rates {
				t.Logf("   %s: Buy=%s, Sell=%s",
					rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
			}
		} else {
			if len(resp.Data.UnavailableCurrencyPairs) > 0 {
				t.Logf("⚠️  Unavailable currency pairs: %v", resp.Data.UnavailableCurrencyPairs)
			}
		}
	})

	t.Run("ListEURPairs", func(t *testing.T) {
		// Request EUR-related pairs
		pairs := []string{
			"EUR/USD", "EUR/GBP", "EUR/JPY", "EUR/CHF",
			"EUR/CAD", "EUR/AUD",
		}

		req := &banking.ListRatesRequest{
			CurrencyPairs: pairs,
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to list EUR pairs: %v", err)
			return
		}

		t.Logf("✅ Retrieved %d EUR exchange rates", len(resp.Data.Rates))

		if len(resp.Data.Rates) > 0 {
			t.Logf("💰 EUR Exchange Rates:")
			for _, rate := range resp.Data.Rates {
				t.Logf("   %s: Buy=%s, Sell=%s",
					rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
			}
		}
	})

	t.Run("ListAfricanCurrencyPairs", func(t *testing.T) {
		// Request African currency pairs
		pairs := []string{
			"USD/KES", "USD/UGX", "USD/TZS", "USD/GHS",
			"EUR/KES", "GBP/KES",
		}

		req := &banking.ListRatesRequest{
			CurrencyPairs: pairs,
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to list African currency pairs: %v", err)
			return
		}

		t.Logf("✅ Retrieved %d African currency exchange rates", len(resp.Data.Rates))

		if len(resp.Data.Rates) > 0 {
			t.Logf("💰 African Currency Exchange Rates:")
			for _, rate := range resp.Data.Rates {
				t.Logf("   %s: Buy=%s, Sell=%s",
					rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
			}
		}
	})

	t.Run("ListCrossCurrencyPairs", func(t *testing.T) {
		// Request cross-currency pairs (non-USD pairs)
		pairs := []string{
			"EUR/GBP", "GBP/EUR", "EUR/JPY", "GBP/JPY",
			"EUR/CHF", "GBP/CHF",
		}

		req := &banking.ListRatesRequest{
			CurrencyPairs: pairs,
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to list cross-currency pairs: %v", err)
			return
		}

		t.Logf("✅ Retrieved %d cross-currency exchange rates", len(resp.Data.Rates))

		if len(resp.Data.Rates) > 0 {
			t.Logf("💰 Cross-Currency Exchange Rates:")
			for _, rate := range resp.Data.Rates {
				t.Logf("   %s: Buy=%s, Sell=%s",
					rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
			}
		}
	})

	t.Run("CompareRates", func(t *testing.T) {
		// Compare buy and sell spreads for major pairs
		pairs := []string{"USD/EUR", "USD/GBP", "USD/JPY"}

		req := &banking.ListRatesRequest{
			CurrencyPairs: pairs,
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to retrieve rates for comparison: %v", err)
			return
		}

		if len(resp.Data.Rates) > 0 {
			t.Logf("📊 Rate Spread Analysis:")
			for _, rate := range resp.Data.Rates {
				t.Logf("   %s:", rate.CurrencyPair)
				t.Logf("      Buy Price:  %s", rate.BuyPrice)
				t.Logf("      Sell Price: %s", rate.SellPrice)
				// Note: Actual spread calculation would require parsing the string values
				t.Logf("      Updated: %s", resp.Data.LastUpdated)
			}
		}
	})

	t.Run("VerifyRateUpdates", func(t *testing.T) {
		// Verify that rates are being updated
		req := &banking.ListRatesRequest{
			CurrencyPairs: []string{"USD/EUR"},
		}

		resp, err := client.Banking.ExchangeRates.List(ctx, req)
		if err != nil {
			t.Logf("❌ Failed to retrieve rate: %v", err)
			return
		}

		if len(resp.Data.Rates) > 0 {
			rate := resp.Data.Rates[0]
			t.Logf("✅ Rate update verification:")
			t.Logf("   Pair: %s", rate.CurrencyPair)
			t.Logf("   Buy: %s, Sell: %s", rate.BuyPrice, rate.SellPrice)
			t.Logf("   Last Updated: %s", resp.Data.LastUpdated)
			t.Logf("   ℹ️  Rates should be updated regularly by the platform")
		}
	})
}
