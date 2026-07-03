package stocktrackingmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockTrackingItem struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	ProductName string    `gorm:"not null" json:"productName"`
	Quantity    int       `gorm:"not null;default:0" json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (s *StockTrackingItem) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}
