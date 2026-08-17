package shipping

import (
	"os"
	"strconv"
)

// Rates is the whole shipping policy: one flat charge, optionally waived above a threshold.
// FreeOverCents == 0 means free shipping is disabled, not that everything ships free.
type Rates struct {
	FlatCents     int64
	FreeOverCents int64
}

func Quote(rates Rates, subtotalCents int64) int64 {
	if subtotalCents <= 0 {
		return 0
	}
	if rates.FreeOverCents > 0 && subtotalCents >= rates.FreeOverCents {
		return 0
	}
	return rates.FlatCents
}

func Total(rates Rates, subtotalCents int64) int64 {
	return subtotalCents + Quote(rates, subtotalCents)
}

func centsFromEnv(key string) int64 {
	parsed, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil || parsed < 0 {
		return 0
	}
	return parsed
}

func RatesFromEnv() Rates {
	return Rates{
		FlatCents:     centsFromEnv("SHIPPING_FLAT_CENTS"),
		FreeOverCents: centsFromEnv("SHIPPING_FREE_OVER_CENTS"),
	}
}
