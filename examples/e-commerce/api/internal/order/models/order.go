package ordermodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            string      `gorm:"primaryKey;type:uuid" json:"id"`
	UserID        *string     `gorm:"type:uuid;index" json:"userId,omitempty"`
	GuestToken    string      `gorm:"index" json:"-"`
	Email         string      `gorm:"not null" json:"email"`
	Status        string      `gorm:"not null;default:pending" json:"status"`
	PaymentMethod string      `gorm:"not null" json:"paymentMethod"`
	PaymentID     string      `json:"-"`
	TotalCents    int64       `gorm:"not null" json:"totalCents"`
	TrackingNo    string      `json:"trackingNo,omitempty"`
	Items         []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"-"`
}

func (o *Order) BeforeCreate(_ *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return nil
}

type OrderItem struct {
	ID        string `gorm:"primaryKey;type:uuid" json:"id"`
	OrderID   string `gorm:"index;not null" json:"orderId"`
	ProductID string `gorm:"type:uuid;not null" json:"productId"`
	TitleTr   string `json:"titleTr"`
	UnitCents int64  `json:"unitCents"`
	Quantity  int    `json:"quantity"`
}

func (i *OrderItem) BeforeCreate(_ *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}
