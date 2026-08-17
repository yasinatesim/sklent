package ordernotifymodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Rule is one "mail this address when an order arrives from that channel" row.
type Rule struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Source    string    `gorm:"not null;index:idx_order_notification_rules_source_recipient,unique" json:"source"`
	Recipient string    `gorm:"not null;index:idx_order_notification_rules_source_recipient,unique" json:"recipient"`
	Enabled   bool      `gorm:"not null;default:true" json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (r *Rule) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (Rule) TableName() string { return "order_notification_rules" }
