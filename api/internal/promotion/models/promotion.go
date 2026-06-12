package promotionmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Promotion struct {
	ID            string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	DiscountType  string    `gorm:"not null" json:"discountType"`
	DiscountValue int       `gorm:"not null" json:"discountValue"`
	ScopeType     string    `gorm:"not null;default:all" json:"scopeType"`
	ProductIDs    string    `json:"-"`
	CategoryIDs   string    `json:"-"`
	MinCartCents  int64     `json:"minCartCents"`
	Active        bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (p *Promotion) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

type Coupon struct {
	ID            string    `gorm:"primaryKey;type:uuid" json:"id"`
	Code          string    `gorm:"uniqueIndex;not null" json:"code"`
	DiscountType  string    `gorm:"not null" json:"discountType"`
	DiscountValue int       `gorm:"not null" json:"discountValue"`
	ScopeType     string    `gorm:"not null;default:all" json:"scopeType"`
	MinCartCents  int64     `json:"minCartCents"`
	Active        bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (c *Coupon) BeforeCreate(_ *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
