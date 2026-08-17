package gib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/invoice/gib"
)

func TestConvertPriceToText(t *testing.T) {
	cases := []struct {
		name  string
		price float64
		want  string
	}{
		{"zero", 0, "YALNIZ SıfırTürkLirası SIFIR KR'DIR."},
		{"round hundred", 100, "YALNIZ YüzTürkLirası SIFIR KR'DIR."},
		{"with kurus", 1234.56, "YALNIZ BinİkiYüzOtuzDörtTürkLirası ElliAltıKR'DIR."},
		{"exactly one thousand says BİN, not BİRBİN", 1000, "YALNIZ BinTürkLirası SIFIR KR'DIR."},
		{"two thousand", 2000, "YALNIZ İkiBinTürkLirası SIFIR KR'DIR."},
		{"millions", 1_000_000, "YALNIZ BirMilyonTürkLirası SIFIR KR'DIR."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, gib.ConvertPriceToText(tc.price))
		})
	}
}

func TestConvertPriceToTextRoundsKurusHalfUp(t *testing.T) {
	// float arithmetic must not drop the last kurus on the legally-printed amount
	assert.Contains(t, gib.ConvertPriceToText(0.99), "DoksanDokuz")
	assert.Contains(t, gib.ConvertPriceToText(19.90), "Doksan")
}
