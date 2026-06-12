package iyzico

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderStore interface {
	GetTotalCents(ctx context.Context, orderID string) (int64, error)
	MarkPaid(ctx context.Context, orderID, paymentID string) error
}

type ReservationStore interface {
	CommitByOrder(ctx context.Context, orderID string) error
	ReleaseByOrder(ctx context.Context, orderID string) error
}

type CallbackResult struct {
	Status    string
	PaidPrice string
	PaymentID string
}

type Verifier interface {
	Verify(ctx context.Context, conversationID, paymentID string) (CallbackResult, error)
}

type Handler struct {
	orders       OrderStore
	reservations ReservationStore
	verifier     Verifier
	frontendBase string
}

func NewHandler(orders OrderStore, reservations ReservationStore, verifier Verifier, frontendBase string) *Handler {
	return &Handler{orders: orders, reservations: reservations, verifier: verifier, frontendBase: frontendBase}
}

// Callback verifies fields, 3DS, order, and exact amount in order; any failure releases + redirects.
func (h *Handler) Callback(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.PostForm("conversationId")
	paymentID := c.PostForm("paymentId")
	mdStatus := c.PostForm("mdStatus")

	if orderID == "" || paymentID == "" {
		h.redirectFE(c, "/odeme/hata")
		return
	}
	if !ThreeDSAuthorized(mdStatus) {
		_ = h.reservations.ReleaseByOrder(ctx, orderID)
		h.redirectFE(c, "/odeme/hata")
		return
	}

	fin, err := h.verifier.Verify(ctx, orderID, paymentID)
	if err != nil {
		_ = h.reservations.ReleaseByOrder(ctx, orderID)
		h.redirectFE(c, "/odeme/hata")
		return
	}

	total, err := h.orders.GetTotalCents(ctx, orderID)
	if err != nil || !AmountMatches(fin.PaidPrice, total) {
		_ = h.reservations.ReleaseByOrder(ctx, orderID)
		h.redirectFE(c, "/odeme/hata")
		return
	}

	if err := h.orders.MarkPaid(ctx, orderID, fin.PaymentID); err != nil {
		h.redirectFE(c, "/odeme/hata")
		return
	}
	if err := h.reservations.CommitByOrder(ctx, orderID); err != nil {
		h.redirectFE(c, "/odeme/hata")
		return
	}
	h.redirectFE(c, "/odeme/basarili?order="+orderID)
}

func (h *Handler) redirectFE(c *gin.Context, path string) {
	c.Redirect(http.StatusFound, h.frontendBase+path)
}
