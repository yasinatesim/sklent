// Package gibmodels defines the wire-level types for GIB's e-Arsiv portal dispatch API.
package gibmodels

import "encoding/json"

// InvoiceItem is one line of a draft invoice sent to GIB.
type InvoiceItem struct {
	Name           string
	Quantity       float64
	UnitType       string
	UnitPrice      float64
	Price          float64
	DiscountRate   float64
	DiscountAmount float64
	DiscountReason string
	Discount       string
	VATRate        float64
	VATAmount      float64
	VATAmountOfTax float64
	TaxRate        float64
}

// InvoiceDetails is the full payload GIB needs to create a draft invoice (RG_BASITFATURA / basit fatura).
type InvoiceDetails struct {
	UUID           string
	DocumentNumber string
	Date           string // DD/MM/YYYY
	Time           string // HH:MM:SS
	InvoiceType    string // default SATIS
	HangiTip       string // default 5000/30000

	Currency     string
	CurrencyRate string

	TaxIDOrTRID string // 11-digit TCKN or 10-digit VKN
	Title       string // corporate unvan
	Name        string
	Surname     string
	FullAddress string
	City        string
	District    string
	Country     string
	ZipCode     string
	TaxOffice   string

	GrandTotal        float64
	TotalDiscount     float64
	TotalVAT          float64
	GrandTotalInclVAT float64
	PaymentTotal      float64

	Items []InvoiceItem
}

// DraftInvoice identifies a just-created draft before GIB's own ETTN can be resolved by listing.
type DraftInvoice struct {
	Date string
	UUID string
}

// InvoiceListItem is one row from GIB's taslaklari-getir (draft/invoice list) response.
type InvoiceListItem struct {
	ETTN               string `json:"ettn"`
	BelgeTarihi        string `json:"belgeTarihi,omitempty"`
	BelgeNumarasi      string `json:"belgeNumarasi,omitempty"`
	AliciUnvanAdSoyad  string `json:"aliciUnvanAdSoyad,omitempty"`
	AliciVknTckn       string `json:"aliciVknTckn,omitempty"`
	SaticiVknTckn      string `json:"saticiVknTckn,omitempty"`
	SaticiUnvanAdSoyad string `json:"saticiUnvanAdSoyad,omitempty"`
	BelgeTuru          string `json:"belgeTuru,omitempty"`
	OnayDurumu         string `json:"onayDurumu,omitempty"`
	FaturaOid          string `json:"faturaOid,omitempty"`
	ToplamTutar        string `json:"toplamTutar,omitempty"`
	TalepDurumColumn   string `json:"talepDurumColumn,omitempty"`
	IptalItiraz        string `json:"iptalItiraz,omitempty"`
	TalepDurum         string `json:"talepDurum,omitempty"`
}

// APIResponse is GIB's generic dispatch envelope. Error is "0" or absent on success.
type APIResponse struct {
	Data     RawJSON      `json:"data"`
	OID      string       `json:"oid,omitempty"`
	Token    string       `json:"token,omitempty"`
	Error    string       `json:"error,omitempty"`
	Messages []GIBMessage `json:"messages,omitempty"`
}

// GIBMessage mirrors GIB's inconsistent message shape — sometimes a plain string, sometimes {type,text}.
type GIBMessage struct {
	Type string
	Text string
}

func (m *GIBMessage) UnmarshalJSON(b []byte) error {
	var asString string
	if err := json.Unmarshal(b, &asString); err == nil {
		m.Text = asString
		return nil
	}
	var asObject struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(b, &asObject); err != nil {
		return err
	}
	m.Type = asObject.Type
	m.Text = asObject.Text
	return nil
}

// RawJSON defers decoding of the polymorphic "data" field until the caller knows its shape.
type RawJSON []byte

func (r *RawJSON) UnmarshalJSON(b []byte) error {
	*r = append((*r)[0:0], b...)
	return nil
}

func (r RawJSON) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return []byte("null"), nil
	}
	return r, nil
}
