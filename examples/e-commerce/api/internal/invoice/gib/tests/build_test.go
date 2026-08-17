package gib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/invoice/gib"
)

func lines() []gib.OrderLine {
	return []gib.OrderLine{
		{Name: "Kolye", Quantity: 2, UnitCents: 10000},
		{Name: "Bileklik", Quantity: 1, UnitCents: 5000},
	}
}

func input() gib.BuildInput {
	return gib.BuildInput{
		DocumentNumber: "VC2026000001",
		OrderStatus:    constants.ORDER_STATUS_PAID,
		Buyer:          gib.Buyer{TaxIDOrTRID: "11111111111", Name: "Ada", Surname: "Lovelace", City: "İstanbul"},
		Lines:          lines(),
		TotalCents:     25000,
		VATRate:        20,
	}
}

func TestBuildRefusesAnUnpayableOrder(t *testing.T) {
	in := input()
	in.OrderStatus = constants.ORDER_STATUS_PENDING

	_, err := gib.BuildInvoiceDetails(in)

	assert.ErrorIs(t, err, gib.ErrOrderNotInvoiceable)
}

func TestBuildRefusesWhenLinesDoNotSumToTheOrderTotal(t *testing.T) {
	in := input()
	in.TotalCents = 999

	_, err := gib.BuildInvoiceDetails(in)

	assert.ErrorIs(t, err, gib.ErrTotalMismatch)
}

func TestBuildRefusesAnInvalidTaxID(t *testing.T) {
	in := input()
	in.Buyer.TaxIDOrTRID = "123"

	_, err := gib.BuildInvoiceDetails(in)

	assert.ErrorIs(t, err, gib.ErrInvalidTaxID)
}

func TestBuildSplitsVATOutOfTheGrossTotal(t *testing.T) {
	details, err := gib.BuildInvoiceDetails(input())

	require.NoError(t, err)
	// 250.00 gross at 20% VAT -> 208.33 net + 41.67 VAT
	assert.InDelta(t, 208.33, details.GrandTotal, 0.01)
	assert.InDelta(t, 41.67, details.TotalVAT, 0.01)
	assert.InDelta(t, 250.00, details.GrandTotalInclVAT, 0.01)
	assert.InDelta(t, 250.00, details.PaymentTotal, 0.01)
}

func TestBuildCarriesEveryOrderLine(t *testing.T) {
	details, err := gib.BuildInvoiceDetails(input())

	require.NoError(t, err)
	require.Len(t, details.Items, 2)
	assert.Equal(t, "Kolye", details.Items[0].Name)
	assert.InDelta(t, 2, details.Items[0].Quantity, 0.001)
	assert.InDelta(t, 20, details.Items[0].VATRate, 0.001)
}

func TestBuildDefaultsTheDocumentShapeGIBExpects(t *testing.T) {
	details, err := gib.BuildInvoiceDetails(input())

	require.NoError(t, err)
	assert.Equal(t, "SATIS", details.InvoiceType)
	assert.Equal(t, "TRY", details.Currency)
	assert.NotEmpty(t, details.UUID)
	assert.Regexp(t, `^\d{2}/\d{2}/\d{4}$`, details.Date)
	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, details.Time)
}
