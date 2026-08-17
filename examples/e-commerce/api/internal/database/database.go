package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yasinatesim/vela-commerce/api/internal/category/models"
	"github.com/yasinatesim/vela-commerce/api/internal/order/models"
	"github.com/yasinatesim/vela-commerce/api/internal/ordernotify/models"
	"github.com/yasinatesim/vela-commerce/api/internal/product/models"
	"github.com/yasinatesim/vela-commerce/api/internal/promotion/models"
	"github.com/yasinatesim/vela-commerce/api/internal/reservation/models"
	"github.com/yasinatesim/vela-commerce/api/internal/returnreq/models"
	"github.com/yasinatesim/vela-commerce/api/internal/review/models"
	"github.com/yasinatesim/vela-commerce/api/internal/stocktracking/models"
	"github.com/yasinatesim/vela-commerce/api/internal/user/models"
)

func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func AutoMigrate(db *gorm.DB) error {
	// Turkish product search compares with unaccent(), so the extension must exist before any
	// query runs; without it the storefront search fails at runtime, not at boot.
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS unaccent").Error; err != nil {
		return err
	}
	return db.AutoMigrate(
		&usermodels.User{},
		&usermodels.RefreshToken{},
		&usermodels.PasswordResetToken{},
		&categorymodels.Category{},
		&productmodels.Product{},
		&ordermodels.Order{},
		&ordermodels.OrderItem{},
		&reservationmodels.Reservation{},
		&promotionmodels.Promotion{},
		&promotionmodels.Coupon{},
		&stocktrackingmodels.StockTrackingItem{},
		&reviewmodels.Review{},
		&ordernotifymodels.Rule{},
		&returnreqmodels.Return{},
	)
}
