package reservation

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/reservation/models"
)

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// Reserve holds stock for an order for RESERVATION_TTL_MINUTES so two buyers cannot take the last unit.
func (s *Service) Reserve(ctx context.Context, orderID, productID string, qty int) error {
	r := reservationmodels.Reservation{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  qty,
		ExpiresAt: time.Now().Add(constants.RESERVATION_TTL_MINUTES * time.Minute),
	}
	return s.db.WithContext(ctx).Create(&r).Error
}

func (s *Service) CommitByOrder(ctx context.Context, orderID string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&reservationmodels.Reservation{}).
		Where("order_id = ? AND committed_at IS NULL AND released_at IS NULL", orderID).
		Update("committed_at", &now).Error
}

func (s *Service) ReleaseByOrder(ctx context.Context, orderID string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&reservationmodels.Reservation{}).
		Where("order_id = ? AND committed_at IS NULL AND released_at IS NULL", orderID).
		Update("released_at", &now).Error
}
