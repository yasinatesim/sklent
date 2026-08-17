package ordernotify

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
)

// hbOrderPush mirrors Hepsiburada's "Create Order" webhook body; HB pushes one envelope per order with its line items.
type hbOrderPush struct {
	Items []struct {
		OrderNumber  json.Number `json:"orderNumber"`
		CustomerName string      `json:"customerName"`
		Quantity     int         `json:"quantity"`
		TotalPrice   struct {
			Amount float64 `json:"amount"`
		} `json:"totalPrice"`
	} `json:"items"`
}

// tyPackagePush mirrors Trendyol's shipment-package webhook body; TY posts the getShipmentPackage shape, so a single package can arrive bare or inside the paginated envelope.
type tyPackagePush struct {
	Content []tyPackage `json:"content"`
}

type tyPackage struct {
	ShipmentPackageID int64   `json:"shipmentPackageId"`
	OrderNumber       string  `json:"orderNumber"`
	CustomerFirstName string  `json:"customerFirstName"`
	CustomerLastName  string  `json:"customerLastName"`
	PackageTotalPrice float64 `json:"packageTotalPrice"`
}

// WebhookHandler serves the order-push contracts the marketplaces call. The example persists nothing: it proves the contract, the basic-auth gate and the alert fan-out.
type WebhookHandler struct{ svc *Service }

func NewWebhookHandler(svc *Service) *WebhookHandler { return &WebhookHandler{svc: svc} }

func centsFromAmount(amount float64) int64 { return int64(amount*100 + 0.5) }

// HBCreateOrder handles POST /webhooks/hb/orders.
func (h *WebhookHandler) HBCreateOrder(c *gin.Context) {
	var push hbOrderPush
	if err := c.ShouldBindJSON(&push); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	for _, order := range groupHBItems(push) {
		h.svc.NotifyNewOrderAsync(constants.ORDER_SOURCE_HB, order.number, order.customer, order.totalCents)
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

type hbOrder struct {
	number     string
	customer   string
	totalCents int64
}

// groupHBItems folds HB's flat item list into one alert per order number.
func groupHBItems(push hbOrderPush) []hbOrder {
	index := map[string]int{}
	orders := []hbOrder{}
	for _, item := range push.Items {
		number := item.OrderNumber.String()
		if number == "" {
			continue
		}
		at, seen := index[number]
		if !seen {
			index[number] = len(orders)
			orders = append(orders, hbOrder{number: number, customer: item.CustomerName})
			at = len(orders) - 1
		}
		orders[at].totalCents += centsFromAmount(item.TotalPrice.Amount)
	}
	return orders
}

// TYOrder handles POST /webhooks/ty/orders.
func (h *WebhookHandler) TYOrder(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	packages, err := decodeTYPackages(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	for _, pkg := range packages {
		number := pkg.OrderNumber
		if number == "" {
			number = strconv.FormatInt(pkg.ShipmentPackageID, 10)
		}
		customer := strings.TrimSpace(pkg.CustomerFirstName + " " + pkg.CustomerLastName)
		h.svc.NotifyNewOrderAsync(constants.ORDER_SOURCE_TY, number, customer, centsFromAmount(pkg.PackageTotalPrice))
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// decodeTYPackages accepts every shape TY is documented to post: the paginated envelope, a bare array, or one package.
func decodeTYPackages(body []byte) ([]tyPackage, error) {
	var envelope tyPackagePush
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Content) > 0 {
		return envelope.Content, nil
	}
	var list []tyPackage
	if err := json.Unmarshal(body, &list); err == nil && len(list) > 0 {
		return list, nil
	}
	var single tyPackage
	if err := json.Unmarshal(body, &single); err != nil {
		return nil, err
	}
	if single.ShipmentPackageID == 0 && single.OrderNumber == "" {
		return nil, nil
	}
	return []tyPackage{single}, nil
}

// Acknowledge answers the status pushes this example stores nothing for, so the marketplace retry loop closes.
func (h *WebhookHandler) Acknowledge(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
