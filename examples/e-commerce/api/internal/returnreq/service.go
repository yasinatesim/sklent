package returnreq

import (
	"context"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	returnreqmodels "github.com/yasinatesim/vela-commerce/api/internal/returnreq/models"
)

type Store interface {
	Create(ctx context.Context, ret returnreqmodels.Return) (returnreqmodels.Return, error)
	OpenForOrder(ctx context.Context, orderID string) (*returnreqmodels.Return, error)
	GetByID(ctx context.Context, id string) (*returnreqmodels.Return, error)
	UpdateStatus(ctx context.Context, id, status, adminNote string) error
}

type OrderLookup interface {
	StatusAndEmail(ctx context.Context, orderID string) (string, string, error)
}

type Service struct {
	store  Store
	orders OrderLookup
}

func NewService(store Store, orders OrderLookup) *Service {
	return &Service{store: store, orders: orders}
}

var validReasons = map[string]struct{}{
	constants.RETURN_REASON_DAMAGED:      {},
	constants.RETURN_REASON_WRONG_ITEM:   {},
	constants.RETURN_REASON_NOT_AS_SHOWN: {},
	constants.RETURN_REASON_CHANGED_MIND: {},
	constants.RETURN_REASON_OTHER:        {},
}

// Request opens a return after checking the order is eligible, the reason is known and no other
// return is already in flight. The email comes from the order, never from the caller.
func (s *Service) Request(ctx context.Context, orderID, reason, comment string) (returnreqmodels.Return, error) {
	if _, ok := validReasons[reason]; !ok {
		return returnreqmodels.Return{}, returnreqmodels.ErrInvalidReason
	}

	orderStatus, email, err := s.orders.StatusAndEmail(ctx, orderID)
	if err != nil {
		return returnreqmodels.Return{}, err
	}
	if !EligibleForReturn(orderStatus) {
		return returnreqmodels.Return{}, returnreqmodels.ErrOrderNotEligible
	}

	open, err := s.store.OpenForOrder(ctx, orderID)
	if err != nil {
		return returnreqmodels.Return{}, err
	}
	if open != nil {
		return returnreqmodels.Return{}, returnreqmodels.ErrReturnExists
	}

	return s.store.Create(ctx, returnreqmodels.Return{
		OrderID: orderID,
		Email:   email,
		Reason:  reason,
		Comment: comment,
		Status:  constants.RETURN_STATUS_REQUESTED,
	})
}

func (s *Service) UpdateStatus(ctx context.Context, id, status, adminNote string) error {
	current, err := s.store.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !CanTransition(current.Status, status) {
		return returnreqmodels.ErrInvalidTransition
	}
	return s.store.UpdateStatus(ctx, id, status, adminNote)
}
