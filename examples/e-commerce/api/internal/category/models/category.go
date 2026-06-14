package categorymodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	NameTr    string    `gorm:"not null" json:"nameTr"`
	NameEn    string    `gorm:"not null" json:"nameEn"`
	Icon      string    `json:"icon"`
	DescTr    string    `json:"descTr"`
	DescEn    string    `json:"descEn"`
	ParentID  *string   `gorm:"index" json:"parentId,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (c *Category) BeforeCreate(_ *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
