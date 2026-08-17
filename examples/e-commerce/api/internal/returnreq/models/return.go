package returnreqmodels

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Return struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	OrderID   string    `gorm:"type:uuid;index;not null" json:"orderId"`
	UserID    *string   `gorm:"type:uuid;index" json:"userId,omitempty"`
	Email     string    `gorm:"not null" json:"email"`
	Reason    string    `gorm:"not null" json:"reason"`
	Comment   string    `json:"comment,omitempty"`
	Status    string    `gorm:"not null;default:requested;index" json:"status"`
	AdminNote string    `json:"adminNote,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (r *Return) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

var (
	ErrReturnNotFound    = errors.New("return not found")
	ErrOrderNotEligible  = errors.New("order is not eligible for a return")
	ErrReturnExists      = errors.New("an open return already exists for this order")
	ErrInvalidReason     = errors.New("invalid return reason")
	ErrInvalidTransition = errors.New("invalid return status transition")
)
