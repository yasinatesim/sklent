package ordernotify

import (
	"context"
	"log/slog"
	"strings"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/ordernotify/models"
)

// SOURCE_LABELS names every channel a rule may target; it doubles as the allowed-source set.
var SOURCE_LABELS = map[string]string{
	constants.ORDER_SOURCE_SITE: "Site",
	constants.ORDER_SOURCE_HB:   "Hepsiburada",
	constants.ORDER_SOURCE_TY:   "Trendyol",
}

// Mailer is the subset of email.Service the notifier needs.
type Mailer interface {
	SendNewOrderAlertAsync(to, sourceLabel, orderNumber, customerName string, totalCents int64)
}

// RuleSource is the subset of Repo the notifier reads, so a test needs no database.
type RuleSource interface {
	ListEnabledBySource(ctx context.Context, source string) ([]ordernotifymodels.Rule, error)
}

type Service struct {
	rules  RuleSource
	mailer Mailer
	log    *slog.Logger
}

func NewService(rules RuleSource, mailer Mailer, log *slog.Logger) *Service {
	return &Service{rules: rules, mailer: mailer, log: log}
}

// IsKnownSource reports whether a rule may target this channel.
func IsKnownSource(source string) bool {
	_, ok := SOURCE_LABELS[source]
	return ok
}

// SourceLabel renders a channel the way the alert email names it.
func SourceLabel(source string) string {
	if label, ok := SOURCE_LABELS[source]; ok {
		return label
	}
	return source
}

// NotifyNewOrderAsync mails every enabled recipient configured for that channel; it never blocks the order path.
func (s *Service) NotifyNewOrderAsync(source, orderNumber, customerName string, totalCents int64) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("new order alert panic", "err", r, "orderNumber", orderNumber)
			}
		}()
		rules, err := s.rules.ListEnabledBySource(context.Background(), source)
		if err != nil {
			s.log.Error("new order alert rules load failed", "err", err, "source", source)
			return
		}
		label := SourceLabel(source)
		sent := map[string]bool{}
		for _, rule := range rules {
			recipient := strings.TrimSpace(rule.Recipient)
			key := strings.ToLower(recipient)
			if recipient == "" || sent[key] {
				continue
			}
			sent[key] = true
			s.mailer.SendNewOrderAlertAsync(recipient, label, orderNumber, customerName, totalCents)
		}
	}()
}
