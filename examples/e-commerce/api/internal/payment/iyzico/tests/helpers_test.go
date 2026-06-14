package iyzico_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/payment/iyzico"
)

func TestAmountMatches(t *testing.T) {
	cases := []struct {
		name       string
		paidPrice  string
		totalCents int64
		want       bool
	}{
		{"exact", "123.45", 12345, true},
		{"one cent under", "123.44", 12345, false},
		{"one cent over", "123.46", 12345, false},
		{"whole number", "100", 10000, true},
		{"garbage", "abc", 10000, false},
		{"empty", "", 10000, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, iyzico.AmountMatches(tc.paidPrice, tc.totalCents))
		})
	}
}

func TestThreeDSAuthorized(t *testing.T) {
	assert.True(t, iyzico.ThreeDSAuthorized("1"))
	assert.False(t, iyzico.ThreeDSAuthorized("0"))
	assert.False(t, iyzico.ThreeDSAuthorized(""))
}
