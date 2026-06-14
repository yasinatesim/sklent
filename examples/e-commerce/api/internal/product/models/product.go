package productmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID            string    `gorm:"primaryKey;type:uuid" json:"id"`
	Slug          string    `gorm:"uniqueIndex;not null" json:"slug"`
	TitleTr       string    `gorm:"not null" json:"titleTr"`
	TitleEn       string    `json:"titleEn"`
	DescriptionTr string    `json:"descriptionTr"`
	DescriptionEn string    `json:"descriptionEn"`
	SeoTitle      string    `json:"seoTitle"`
	PriceCents    int64     `gorm:"not null" json:"priceCents"`
	OldPriceCents int64     `json:"oldPriceCents"`
	Stock         int       `gorm:"not null;default:0" json:"stock"`
	CategoryID    *string   `gorm:"type:uuid;index" json:"categoryId,omitempty"`
	CategorySlug  string    `gorm:"index" json:"categorySlug"`
	Badge         string    `json:"badge"`
	Seller        string    `json:"seller"`
	DominantColor string    `json:"dominantColor"`
	Material      string    `json:"material"`
	Published     bool      `gorm:"not null;default:false" json:"published"`
	ImageURL      string    `json:"imageUrl"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (p *Product) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}
