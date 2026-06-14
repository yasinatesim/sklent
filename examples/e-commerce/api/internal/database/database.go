package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yasinatesim/vela-commerce/api/internal/category/models"
	"github.com/yasinatesim/vela-commerce/api/internal/order/models"
	"github.com/yasinatesim/vela-commerce/api/internal/product/models"
	"github.com/yasinatesim/vela-commerce/api/internal/promotion/models"
	"github.com/yasinatesim/vela-commerce/api/internal/reservation/models"
	"github.com/yasinatesim/vela-commerce/api/internal/user/models"
)

func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermodels.User{},
		&usermodels.RefreshToken{},
		&categorymodels.Category{},
		&productmodels.Product{},
		&ordermodels.Order{},
		&ordermodels.OrderItem{},
		&reservationmodels.Reservation{},
		&promotionmodels.Promotion{},
		&promotionmodels.Coupon{},
	)
}
