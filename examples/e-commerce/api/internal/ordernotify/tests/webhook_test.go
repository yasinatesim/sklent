package ordernotify_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/ordernotify"
	"github.com/yasinatesim/vela-commerce/api/internal/ordernotify/models"
)

const (
	webhookUser     = "mp-push"
	webhookPassword = "s3cret"
)

type sentAlert struct {
	to          string
	sourceLabel string
	orderNumber string
	customer    string
	totalCents  int64
}

type fakeMailer struct {
	mu    sync.Mutex
	sends []sentAlert
}

func (f *fakeMailer) SendNewOrderAlertAsync(to, sourceLabel, orderNumber, customerName string, totalCents int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sends = append(f.sends, sentAlert{to, sourceLabel, orderNumber, customerName, totalCents})
}

func (f *fakeMailer) received() []sentAlert {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]sentAlert(nil), f.sends...)
}

type fakeRules struct {
	bySource map[string][]ordernotifymodels.Rule
	err      error
}

func (f *fakeRules) ListEnabledBySource(_ context.Context, source string) ([]ordernotifymodels.Rule, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.bySource[source], nil
}

func newWebhookEngine(rules *fakeRules) (*gin.Engine, *fakeMailer) {
	gin.SetMode(gin.TestMode)
	mailer := &fakeMailer{}
	svc := ordernotify.NewService(rules, mailer, slog.New(slog.NewTextHandler(io.Discard, nil)))
	h := ordernotify.NewWebhookHandler(svc)

	r := gin.New()
	basicAuth := gin.BasicAuth(gin.Accounts{webhookUser: webhookPassword})
	hb := r.Group("/webhooks/hb", basicAuth)
	hb.POST("/orders", h.HBCreateOrder)
	hb.PUT("/packages/:packageNumber/deliver", h.Acknowledge)
	ty := r.Group("/webhooks/ty", basicAuth)
	ty.POST("/orders", h.TYOrder)
	return r, mailer
}

func enabledRule(source, recipient string) ordernotifymodels.Rule {
	return ordernotifymodels.Rule{Source: source, Recipient: recipient, Enabled: true}
}

func post(r *gin.Engine, path, body string, authenticated bool) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		req.SetBasicAuth(webhookUser, webhookPassword)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// waitForAlerts polls because the notifier deliberately runs off the request path.
func waitForAlerts(t *testing.T, mailer *fakeMailer, want int) []sentAlert {
	t.Helper()
	for range 200 {
		if got := mailer.received(); len(got) >= want {
			return got
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("expected %d alert(s), got %d", want, len(mailer.received()))
	return nil
}

func TestHBCreateOrder_AlertsOnlyTheRecipientsConfiguredForHB(t *testing.T) {
	r, mailer := newWebhookEngine(&fakeRules{bySource: map[string][]ordernotifymodels.Rule{
		constants.ORDER_SOURCE_HB: {enabledRule(constants.ORDER_SOURCE_HB, "ops@vela.test")},
		constants.ORDER_SOURCE_TY: {enabledRule(constants.ORDER_SOURCE_TY, "other@vela.test")},
	}})

	body := `{"items":[{"orderNumber":"1451124210","customerName":"Ahmet Aslan","quantity":1,"totalPrice":{"amount":14.95}},
{"orderNumber":"1451124210","customerName":"Ahmet Aslan","quantity":2,"totalPrice":{"amount":10}}]}`
	w := post(r, "/webhooks/hb/orders", body, true)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}
	alerts := waitForAlerts(t, mailer, 1)
	if len(alerts) != 1 {
		t.Fatalf("one order number must produce one alert, got %d", len(alerts))
	}
	if alerts[0].to != "ops@vela.test" {
		t.Errorf("only the HB rule applies, got %q", alerts[0].to)
	}
	if alerts[0].sourceLabel != "Hepsiburada" {
		t.Errorf("want the Hepsiburada label, got %q", alerts[0].sourceLabel)
	}
	if alerts[0].totalCents != 2495 {
		t.Errorf("line totals of one order must add up, got %d", alerts[0].totalCents)
	}
}

func TestHBCreateOrder_RejectsUnauthenticatedPush(t *testing.T) {
	r, mailer := newWebhookEngine(&fakeRules{bySource: map[string][]ordernotifymodels.Rule{
		constants.ORDER_SOURCE_HB: {enabledRule(constants.ORDER_SOURCE_HB, "ops@vela.test")},
	}})

	w := post(r, "/webhooks/hb/orders", `{"items":[]}`, false)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
	if len(mailer.received()) != 0 {
		t.Error("an unauthenticated push must not raise an alert")
	}
}

func TestHBCreateOrder_RejectsMalformedBody(t *testing.T) {
	r, _ := newWebhookEngine(&fakeRules{})

	if w := post(r, "/webhooks/hb/orders", "not json", true); w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestTYOrder_AcceptsEveryDocumentedPayloadShape(t *testing.T) {
	pkg := `{"shipmentPackageId":4242,"orderNumber":"TY-99","customerFirstName":"Ahmet","customerLastName":"Aslan","packageTotalPrice":149.9}`
	cases := map[string]string{
		"paginated envelope": `{"content":[` + pkg + `]}`,
		"bare array":         `[` + pkg + `]`,
		"single package":     pkg,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			r, mailer := newWebhookEngine(&fakeRules{bySource: map[string][]ordernotifymodels.Rule{
				constants.ORDER_SOURCE_TY: {enabledRule(constants.ORDER_SOURCE_TY, "ops@vela.test")},
			}})

			w := post(r, "/webhooks/ty/orders", body, true)

			if w.Code != http.StatusOK {
				t.Fatalf("want 200, got %d", w.Code)
			}
			alerts := waitForAlerts(t, mailer, 1)
			if alerts[0].orderNumber != "TY-99" || alerts[0].totalCents != 14990 {
				t.Errorf("unexpected alert: %+v", alerts[0])
			}
			if alerts[0].customer != "Ahmet Aslan" {
				t.Errorf("want the joined customer name, got %q", alerts[0].customer)
			}
		})
	}
}

func TestNotifyNewOrderAsync_MailsOneAddressOnceRegardlessOfCase(t *testing.T) {
	rules := &fakeRules{bySource: map[string][]ordernotifymodels.Rule{
		constants.ORDER_SOURCE_SITE: {
			enabledRule(constants.ORDER_SOURCE_SITE, "ops@vela.test"),
			enabledRule(constants.ORDER_SOURCE_SITE, " OPS@vela.test "),
		},
	}}
	mailer := &fakeMailer{}
	svc := ordernotify.NewService(rules, mailer, slog.New(slog.NewTextHandler(io.Discard, nil)))

	svc.NotifyNewOrderAsync(constants.ORDER_SOURCE_SITE, "ord-1", "ops@vela.test", 15000)

	alerts := waitForAlerts(t, mailer, 1)
	time.Sleep(20 * time.Millisecond)
	if len(mailer.received()) != 1 {
		t.Fatalf("the same address twice must send once: %+v", mailer.received())
	}
	if alerts[0].sourceLabel != "Site" {
		t.Errorf("want the Site label, got %q", alerts[0].sourceLabel)
	}
}

func TestNotifyNewOrderAsync_RuleLoadFailureSendsNothing(t *testing.T) {
	mailer := &fakeMailer{}
	svc := ordernotify.NewService(&fakeRules{err: errors.New("db down")}, mailer,
		slog.New(slog.NewTextHandler(io.Discard, nil)))

	svc.NotifyNewOrderAsync(constants.ORDER_SOURCE_SITE, "ord-1", "ops@vela.test", 15000)

	time.Sleep(50 * time.Millisecond)
	if len(mailer.received()) != 0 {
		t.Error("a failed rule lookup must not fabricate an alert")
	}
}

func TestStatusPushIsAcknowledged(t *testing.T) {
	r, _ := newWebhookEngine(&fakeRules{})

	req := httptest.NewRequest(http.MethodPut, "/webhooks/hb/packages/013105889/deliver", strings.NewReader(`{}`))
	req.SetBasicAuth(webhookUser, webhookPassword)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204 so the marketplace stops retrying, got %d", w.Code)
	}
}
