package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/token"
	"github.com/yasinatesim/vela-commerce/api/internal/user/models"
)

var ErrNotFound = errors.New("auth: not found")

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) CreateUser(ctx context.Context, u *usermodels.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (*usermodels.User, error) {
	var u usermodels.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *Repo) FindByID(ctx context.Context, id string) (*usermodels.User, error) {
	var u usermodels.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *Repo) StoreRefresh(ctx context.Context, userID, rawToken string, ttl time.Duration) error {
	rt := usermodels.RefreshToken{
		UserID:    userID,
		TokenHash: token.Hash(rawToken),
		ExpiresAt: time.Now().Add(ttl),
	}
	return r.db.WithContext(ctx).Create(&rt).Error
}

// RotateRefresh revokes the presented token and issues nothing here; the handler stores a new one.
func (r *Repo) RotateRefresh(ctx context.Context, rawToken string) (string, error) {
	var rt usermodels.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND revoked_at IS NULL", token.Hash(rawToken)).First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if time.Now().After(rt.ExpiresAt) {
		return "", ErrNotFound
	}
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&rt).Update("revoked_at", &now).Error; err != nil {
		return "", err
	}
	return rt.UserID, nil
}
