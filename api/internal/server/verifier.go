package server

import (
	"context"
	"fmt"

	"github.com/yasinatesim/vela-commerce/api/internal/payment/iyzico"
)

// sandboxVerifier echoes paymentId (the decimal TL amount) back as paid so 3DS works without a card.
type sandboxVerifier struct{}

func (sandboxVerifier) Verify(_ context.Context, conversationID, paymentID string) (iyzico.CallbackResult, error) {
	return iyzico.CallbackResult{
		Status:    "success",
		PaidPrice: paymentID,
		PaymentID: fmt.Sprintf("sandbox-%s", conversationID),
	}, nil
}
