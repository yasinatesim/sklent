package reviewmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID         string    `gorm:"primaryKey;type:uuid" json:"id"`
	ProductID  string    `gorm:"type:uuid;not null;index" json:"productId"`
	AuthorName string    `gorm:"not null" json:"authorName"`
	Rating     int       `gorm:"not null" json:"rating"`
	Comment    string    `json:"comment"`
	Status     string    `gorm:"not null;default:pending" json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (r *Review) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}
