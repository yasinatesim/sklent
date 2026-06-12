package reservationmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	ID          string     `gorm:"primaryKey;type:uuid" json:"id"`
	OrderID     string     `gorm:"index;not null" json:"orderId"`
	ProductID   string     `gorm:"type:uuid;index;not null" json:"productId"`
	Quantity    int        `gorm:"not null" json:"quantity"`
	ExpiresAt   time.Time  `gorm:"index" json:"expiresAt"`
	CommittedAt *time.Time `json:"committedAt,omitempty"`
	ReleasedAt  *time.Time `json:"releasedAt,omitempty"`
	CreatedAt   time.Time  `json:"-"`
}

func (r *Reservation) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}
