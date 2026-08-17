package shipping_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/shipping"
)

func TestQuoteFlatRate(t *testing.T) {
	rates := shipping.Rates{FlatCents: 2999, FreeOverCents: 50000}

	cases := []struct {
		name          string
		subtotalCents int64
		want          int64
	}{
		{"below the threshold pays the flat rate", 49999, 2999},
		{"exactly at the threshold ships free", 50000, 0},
		{"above the threshold ships free", 50001, 0},
		{"an empty cart is never charged shipping", 0, 0},
		{"a negative subtotal is treated as empty, not as a discount", -100, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, shipping.Quote(rates, tc.subtotalCents))
		})
	}
}

func TestQuoteWithFreeShippingDisabled(t *testing.T) {
	rates := shipping.Rates{FlatCents: 2999, FreeOverCents: 0}

	assert.Equal(t, int64(2999), shipping.Quote(rates, 1))
	assert.Equal(t, int64(2999), shipping.Quote(rates, 10_000_000),
		"a zero threshold means free shipping is off, not that everything ships free")
	assert.Equal(t, int64(0), shipping.Quote(rates, 0))
}

func TestTotalAddsShippingToTheSubtotal(t *testing.T) {
	rates := shipping.Rates{FlatCents: 2999, FreeOverCents: 50000}

	assert.Equal(t, int64(12999), shipping.Total(rates, 10000))
	assert.Equal(t, int64(60000), shipping.Total(rates, 60000))
}

func TestRatesFromEnvFallsBackToSafeDefaults(t *testing.T) {
	t.Setenv("SHIPPING_FLAT_CENTS", "")
	t.Setenv("SHIPPING_FREE_OVER_CENTS", "")
	rates := shipping.RatesFromEnv()
	assert.GreaterOrEqual(t, rates.FlatCents, int64(0))

	t.Setenv("SHIPPING_FLAT_CENTS", "1500")
	t.Setenv("SHIPPING_FREE_OVER_CENTS", "20000")
	rates = shipping.RatesFromEnv()
	assert.Equal(t, int64(1500), rates.FlatCents)
	assert.Equal(t, int64(20000), rates.FreeOverCents)

	t.Setenv("SHIPPING_FLAT_CENTS", "not-a-number")
	rates = shipping.RatesFromEnv()
	assert.Equal(t, int64(0), rates.FlatCents, "garbage must not become a random charge")
}
