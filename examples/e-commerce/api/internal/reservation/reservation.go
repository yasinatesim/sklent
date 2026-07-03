package reservation

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/reservation/models"
)

type ProductStock interface {
	DecrementStock(ctx context.Context, productID string, qty int) (newStock, threshold int, title string, err error)
}

type Mailer interface {
	SendLowStockAlertAsync(adminEmail, productTitle string, currentStock int)
}

type Service struct {
	db         *gorm.DB
	stock      ProductStock
	mailer     Mailer
	adminEmail string
}

func NewService(db *gorm.DB, stock ProductStock, mailer Mailer, adminEmail string) *Service {
	return &Service{db: db, stock: stock, mailer: mailer, adminEmail: adminEmail}
}

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

// CommitByOrder marks reservations committed and decrements real product stock, firing a low-stock alert if needed.
func (s *Service) CommitByOrder(ctx context.Context, orderID string) error {
	var pending []reservationmodels.Reservation
	if err := s.db.WithContext(ctx).
		Where("order_id = ? AND committed_at IS NULL AND released_at IS NULL", orderID).
		Find(&pending).Error; err != nil {
		return err
	}

	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&reservationmodels.Reservation{}).
		Where("order_id = ? AND committed_at IS NULL AND released_at IS NULL", orderID).
		Update("committed_at", &now).Error; err != nil {
		return err
	}

	for _, r := range pending {
		newStock, threshold, title, err := s.stock.DecrementStock(ctx, r.ProductID, r.Quantity)
		if err != nil {
			continue
		}
		if newStock <= threshold {
			s.mailer.SendLowStockAlertAsync(s.adminEmail, title, newStock)
		}
	}
	return nil
}

func (s *Service) ReleaseByOrder(ctx context.Context, orderID string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&reservationmodels.Reservation{}).
		Where("order_id = ? AND committed_at IS NULL AND released_at IS NULL", orderID).
		Update("released_at", &now).Error
}
