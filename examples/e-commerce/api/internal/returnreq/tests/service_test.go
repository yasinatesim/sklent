package returnreq_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/returnreq"
	returnreqmodels "github.com/yasinatesim/vela-commerce/api/internal/returnreq/models"
)

type stubStore struct {
	created returnreqmodels.Return
	open    *returnreqmodels.Return
	byID    *returnreqmodels.Return
	updated struct {
		id, status, note string
		called           bool
	}
}

func (s *stubStore) Create(_ context.Context, ret returnreqmodels.Return) (returnreqmodels.Return, error) {
	s.created = ret
	return ret, nil
}
func (s *stubStore) OpenForOrder(_ context.Context, _ string) (*returnreqmodels.Return, error) {
	return s.open, nil
}
func (s *stubStore) GetByID(_ context.Context, _ string) (*returnreqmodels.Return, error) {
	if s.byID == nil {
		return nil, returnreqmodels.ErrReturnNotFound
	}
	return s.byID, nil
}
func (s *stubStore) UpdateStatus(_ context.Context, id, status, note string) error {
	s.updated.id, s.updated.status, s.updated.note, s.updated.called = id, status, note, true
	return nil
}

type stubOrders struct {
	status string
	email  string
}

func (s stubOrders) StatusAndEmail(_ context.Context, _ string) (string, string, error) {
	return s.status, s.email, nil
}

func newService(store *stubStore, orderStatus string) *returnreq.Service {
	return returnreq.NewService(store, stubOrders{status: orderStatus, email: "buyer@example.com"})
}

func TestRequestRejectsUnshippedOrders(t *testing.T) {
	store := &stubStore{}
	svc := newService(store, constants.ORDER_STATUS_PAID)

	_, err := svc.Request(context.Background(), "ord-1", constants.RETURN_REASON_DAMAGED, "")

	assert.ErrorIs(t, err, returnreqmodels.ErrOrderNotEligible)
	assert.Empty(t, store.created.ID)
}

func TestRequestRejectsAnUnknownReason(t *testing.T) {
	svc := newService(&stubStore{}, constants.ORDER_STATUS_SHIPPED)

	_, err := svc.Request(context.Background(), "ord-1", "because", "")

	assert.ErrorIs(t, err, returnreqmodels.ErrInvalidReason)
}

func TestRequestRefusesASecondOpenReturn(t *testing.T) {
	open := &returnreqmodels.Return{ID: "r-1", Status: constants.RETURN_STATUS_REQUESTED}
	svc := newService(&stubStore{open: open}, constants.ORDER_STATUS_SHIPPED)

	_, err := svc.Request(context.Background(), "ord-1", constants.RETURN_REASON_DAMAGED, "")

	assert.ErrorIs(t, err, returnreqmodels.ErrReturnExists)
}

func TestRequestStoresTheOrdersEmailNotAClientSuppliedOne(t *testing.T) {
	store := &stubStore{}
	svc := newService(store, constants.ORDER_STATUS_SHIPPED)

	got, err := svc.Request(context.Background(), "ord-1", constants.RETURN_REASON_DAMAGED, "box was crushed")

	require.NoError(t, err)
	assert.Equal(t, "buyer@example.com", got.Email)
	assert.Equal(t, constants.RETURN_STATUS_REQUESTED, got.Status)
	assert.Equal(t, "box was crushed", got.Comment)
}

func TestUpdateStatusRefusesAnIllegalTransition(t *testing.T) {
	store := &stubStore{byID: &returnreqmodels.Return{ID: "r-1", Status: constants.RETURN_STATUS_REQUESTED}}
	svc := newService(store, constants.ORDER_STATUS_SHIPPED)

	err := svc.UpdateStatus(context.Background(), "r-1", constants.RETURN_STATUS_REFUNDED, "")

	assert.ErrorIs(t, err, returnreqmodels.ErrInvalidTransition)
	assert.False(t, store.updated.called, "an illegal transition must not reach the store")
}

func TestUpdateStatusAppliesALegalTransition(t *testing.T) {
	store := &stubStore{byID: &returnreqmodels.Return{ID: "r-1", Status: constants.RETURN_STATUS_APPROVED}}
	svc := newService(store, constants.ORDER_STATUS_SHIPPED)

	err := svc.UpdateStatus(context.Background(), "r-1", constants.RETURN_STATUS_REFUNDED, "money sent")

	require.NoError(t, err)
	assert.True(t, store.updated.called)
	assert.Equal(t, constants.RETURN_STATUS_REFUNDED, store.updated.status)
	assert.Equal(t, "money sent", store.updated.note)
}
