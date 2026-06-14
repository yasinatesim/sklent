package usermodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey;type:uuid"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null;default:user"`
	FullName     string
	ClosedAt     *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

type RefreshToken struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	UserID    string `gorm:"index;not null"`
	TokenHash string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (r *RefreshToken) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}
