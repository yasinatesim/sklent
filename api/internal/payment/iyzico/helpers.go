package iyzico

import (
	"math"
	"strconv"
)

// AmountMatches reports whether paidPrice (decimal TL) equals totalCents; a one-cent gap fails.
func AmountMatches(paidPrice string, totalCents int64) bool {
	paid, err := strconv.ParseFloat(paidPrice, 64)
	if err != nil {
		return false
	}
	return math.Abs(paid*100-float64(totalCents)) <= 0.001
}

// ThreeDSAuthorized reports whether the 3DS verification step succeeded.
func ThreeDSAuthorized(mdStatus string) bool {
	return mdStatus == "1"
}
