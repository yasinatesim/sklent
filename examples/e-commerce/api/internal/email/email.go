package email

import (
	"context"
	"fmt"
	"log/slog"
)

type Message struct {
	To      string
	Subject string
	HTML    string
}

type Sender interface {
	Send(ctx context.Context, m Message) error
}

// LogSender is used when no SMTP host is configured: it logs instead of sending.
type LogSender struct{ Log *slog.Logger }

func (s LogSender) Send(_ context.Context, m Message) error {
	s.Log.Info("email (log only)", "to", m.To, "subject", m.Subject)
	return nil
}

type Service struct {
	sender Sender
	log    *slog.Logger
}

func NewService(sender Sender, log *slog.Logger) *Service {
	return &Service{sender: sender, log: log}
}

type OrderSummary struct {
	ID         string
	Email      string
	TotalCents int64
}

// SendOrderConfirmationAsync never blocks the payment path; a panic is recovered and logged.
func (s *Service) SendOrderConfirmationAsync(order OrderSummary) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("order confirmation email panic", "err", r, "orderID", order.ID)
			}
		}()
		m := Message{
			To:      order.Email,
			Subject: "Siparişiniz alındı",
			HTML:    fmt.Sprintf("<p>Sipariş %s tutarı %d kuruş.</p>", order.ID, order.TotalCents),
		}
		if err := s.sender.Send(context.Background(), m); err != nil {
			s.log.Error("order confirmation email send failed", "err", err, "orderID", order.ID)
		}
	}()
}

// SendPasswordResetAsync never blocks the request path; a panic is recovered and logged.
func (s *Service) SendPasswordResetAsync(to, resetURL string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("password reset email panic", "err", r, "to", to)
			}
		}()
		m := Message{
			To:      to,
			Subject: "Şifre sıfırlama isteği",
			HTML:    fmt.Sprintf("<p>Şifrenizi sıfırlamak için <a href=\"%s\">bu bağlantıya</a> tıklayın. Bağlantı 30 dakika geçerlidir.</p>", resetURL),
		}
		if err := s.sender.Send(context.Background(), m); err != nil {
			s.log.Error("password reset email send failed", "err", err, "to", to)
		}
	}()
}

// SendLowStockAlertAsync notifies the admin mailbox when a product's stock crosses its threshold.
func (s *Service) SendLowStockAlertAsync(adminEmail, productTitle string, currentStock int) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("low stock alert email panic", "err", r, "product", productTitle)
			}
		}()
		m := Message{
			To:      adminEmail,
			Subject: "Düşük stok uyarısı",
			HTML:    fmt.Sprintf("<p>%q ürününün stoğu %d adede düştü.</p>", productTitle, currentStock),
		}
		if err := s.sender.Send(context.Background(), m); err != nil {
			s.log.Error("low stock alert email send failed", "err", err, "product", productTitle)
		}
	}()
}
