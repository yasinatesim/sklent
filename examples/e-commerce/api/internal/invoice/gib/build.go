package gib

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	gibmodels "github.com/yasinatesim/vela-commerce/api/internal/invoice/gib/models"
)

var (
	ErrOrderNotInvoiceable = errors.New("order status does not allow invoicing")
	ErrTotalMismatch       = errors.New("invoice lines do not sum to the order total")
	ErrInvalidTaxID        = errors.New("tax id must be an 11-digit TCKN or a 10-digit VKN")
)

type OrderLine struct {
	Name      string
	Quantity  float64
	UnitCents int64
}

type Buyer struct {
	TaxIDOrTRID string
	Title       string
	Name        string
	Surname     string
	FullAddress string
	City        string
	District    string
	TaxOffice   string
}

type BuildInput struct {
	DocumentNumber string
	OrderStatus    string
	Buyer          Buyer
	Lines          []OrderLine
	TotalCents     int64
	VATRate        float64
}

// Money already left the buyer's account before an invoice exists, so only paid or shipped
// orders may be invoiced.
var invoiceableOrderStatuses = map[string]struct{}{
	constants.ORDER_STATUS_PAID:    {},
	constants.ORDER_STATUS_SHIPPED: {},
}

const (
	TCKN_LENGTH     = 11
	VKN_LENGTH      = 10
	DEFAULT_COUNTRY = "Türkiye"
	CURRENCY_TRY    = "TRY"
	INVOICE_TYPE    = "SATIS"
	UNIT_TYPE_PIECE = "C62"
	PERCENT         = 100.0
	CENTS_IN_LIRA   = 100.0
)

func CanIssue(orderStatus string) bool {
	_, ok := invoiceableOrderStatuses[orderStatus]
	return ok
}

func validTaxID(id string) bool {
	if len(id) != TCKN_LENGTH && len(id) != VKN_LENGTH {
		return false
	}
	for _, r := range id {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func linesTotalCents(lines []OrderLine) int64 {
	var total int64
	for _, line := range lines {
		total += int64(line.Quantity) * line.UnitCents
	}
	return total
}

// BuildInvoiceDetails renders an order as the RG_BASITFATURA payload GIB's e-Arsiv portal expects.
// GIB wants VAT split out of a gross total, so the line prices are the tax-exclusive halves.
func BuildInvoiceDetails(in BuildInput) (gibmodels.InvoiceDetails, error) {
	if !CanIssue(in.OrderStatus) {
		return gibmodels.InvoiceDetails{}, ErrOrderNotInvoiceable
	}
	if linesTotalCents(in.Lines) != in.TotalCents {
		return gibmodels.InvoiceDetails{}, ErrTotalMismatch
	}
	if !validTaxID(in.Buyer.TaxIDOrTRID) {
		return gibmodels.InvoiceDetails{}, ErrInvalidTaxID
	}

	gross := float64(in.TotalCents) / CENTS_IN_LIRA
	net := gross / (1 + in.VATRate/PERCENT)
	vat := gross - net

	items := make([]gibmodels.InvoiceItem, 0, len(in.Lines))
	for _, line := range in.Lines {
		lineGross := float64(line.UnitCents) / CENTS_IN_LIRA * line.Quantity
		lineNet := lineGross / (1 + in.VATRate/PERCENT)
		items = append(items, gibmodels.InvoiceItem{
			Name:      line.Name,
			Quantity:  line.Quantity,
			UnitType:  UNIT_TYPE_PIECE,
			UnitPrice: lineNet / line.Quantity,
			Price:     lineNet,
			VATRate:   in.VATRate,
			VATAmount: lineGross - lineNet,
		})
	}

	now := time.Now()
	return gibmodels.InvoiceDetails{
		UUID:              uuid.NewString(),
		DocumentNumber:    in.DocumentNumber,
		Date:              now.Format("02/01/2006"),
		Time:              now.Format("15:04:05"),
		InvoiceType:       INVOICE_TYPE,
		Currency:          CURRENCY_TRY,
		TaxIDOrTRID:       in.Buyer.TaxIDOrTRID,
		Title:             in.Buyer.Title,
		Name:              in.Buyer.Name,
		Surname:           in.Buyer.Surname,
		FullAddress:       in.Buyer.FullAddress,
		City:              in.Buyer.City,
		District:          in.Buyer.District,
		Country:           DEFAULT_COUNTRY,
		TaxOffice:         in.Buyer.TaxOffice,
		GrandTotal:        net,
		TotalVAT:          vat,
		GrandTotalInclVAT: gross,
		PaymentTotal:      gross,
		Items:             items,
	}, nil
}
