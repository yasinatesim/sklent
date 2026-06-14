package invoice_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/invoice"
)

func TestIsGibTarget(t *testing.T) {
	cases := []struct {
		target string
		want   bool
	}{
		{"https://earsivportal.efatura.gov.tr/login", true},
		{"https://earsivportaltest.efatura.gov.tr/x", true},
		{"http://earsivportal.efatura.gov.tr/login", false},
		{"https://evil.example.com/", false},
		{"https://earsivportal.efatura.gov.tr.evil.com/", false},
		{"not-a-url", false},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, invoice.IsGibTarget(tc.target), tc.target)
	}
}
