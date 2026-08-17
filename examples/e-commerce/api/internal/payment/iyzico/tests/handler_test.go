package iyzico_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/payment/iyzico"
)

type stubOrders struct {
	total    int64
	totalErr error
	paidErr  error
	paidWith string
}

func (s *stubOrders) GetTotalCents(_ context.Context, _ string) (int64, error) {
	return s.total, s.totalErr
}

func (s *stubOrders) MarkPaid(_ context.Context, _, paymentID string) error {
	s.paidWith = paymentID
	return s.paidErr
}

type stubReservations struct {
	committed bool
	released  bool
	commitErr error
}

func (s *stubReservations) CommitByOrder(_ context.Context, _ string) error {
	s.committed = true
	return s.commitErr
}

func (s *stubReservations) ReleaseByOrder(_ context.Context, _ string) error {
	s.released = true
	return nil
}

type stubVerifier struct {
	result iyzico.CallbackResult
	err    error
}

func (s stubVerifier) Verify(_ context.Context, _, _ string) (iyzico.CallbackResult, error) {
	return s.result, s.err
}

const frontendBase = "https://shop.example"

func postCallback(t *testing.T, h *iyzico.Handler, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/callback", h.Callback)

	req := httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCallbackRedirectsToLocaleNeutralPaths(t *testing.T) {
	turkish := []string{"odeme", "basarili", "hata"}

	cases := []struct {
		name        string
		form        url.Values
		orders      *stubOrders
		verifier    stubVerifier
		wantPath    string
		wantCommit  bool
		wantRelease bool
	}{
		{
			name:     "missing fields",
			form:     url.Values{"conversationId": {""}, "paymentId": {""}},
			orders:   &stubOrders{},
			wantPath: "/checkout/error",
		},
		{
			name:        "3ds not authorized releases the reservation",
			form:        url.Values{"conversationId": {"ord-1"}, "paymentId": {"pay-1"}, "mdStatus": {"0"}},
			orders:      &stubOrders{},
			wantPath:    "/checkout/error",
			wantRelease: true,
		},
		{
			name:        "amount mismatch releases the reservation",
			form:        url.Values{"conversationId": {"ord-1"}, "paymentId": {"pay-1"}, "mdStatus": {"1"}},
			orders:      &stubOrders{total: 12345},
			verifier:    stubVerifier{result: iyzico.CallbackResult{Status: "success", PaidPrice: "1.00", PaymentID: "pay-1"}},
			wantPath:    "/checkout/error",
			wantRelease: true,
		},
		{
			name:       "authorized and matching amount commits and redirects to success",
			form:       url.Values{"conversationId": {"ord-1"}, "paymentId": {"pay-1"}, "mdStatus": {"1"}},
			orders:     &stubOrders{total: 12345},
			verifier:   stubVerifier{result: iyzico.CallbackResult{Status: "success", PaidPrice: "123.45", PaymentID: "pay-9"}},
			wantPath:   "/checkout/success",
			wantCommit: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := &stubReservations{}
			h := iyzico.NewHandler(tc.orders, res, tc.verifier, frontendBase)

			rec := postCallback(t, h, tc.form)

			require.Equal(t, http.StatusFound, rec.Code)
			location := rec.Header().Get("Location")
			assert.True(t, strings.HasPrefix(location, frontendBase+tc.wantPath),
				"got %q, want prefix %q", location, frontendBase+tc.wantPath)

			for _, seg := range turkish {
				assert.NotContains(t, location, "/"+seg, "redirect paths stay locale-neutral")
			}

			assert.Equal(t, tc.wantCommit, res.committed)
			assert.Equal(t, tc.wantRelease, res.released)
		})
	}
}

func TestCallbackMarksPaidWithTheVerifiedPaymentID(t *testing.T) {
	orders := &stubOrders{total: 12345}
	res := &stubReservations{}
	h := iyzico.NewHandler(orders, res,
		stubVerifier{result: iyzico.CallbackResult{Status: "success", PaidPrice: "123.45", PaymentID: "verified-id"}},
		frontendBase)

	postCallback(t, h, url.Values{"conversationId": {"ord-1"}, "paymentId": {"posted-id"}, "mdStatus": {"1"}})

	assert.Equal(t, "verified-id", orders.paidWith,
		"must trust the verifier's payment id, never the posted one")
}
