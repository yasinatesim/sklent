package promotion

import (
	"context"
	"errors"

	"gorm.io/gorm"

	promotionmodels "github.com/yasinatesim/vela-commerce/api/internal/promotion/models"
)

var ErrNotFound = errors.New("promotion: not found")

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListPromotions(ctx context.Context) ([]promotionmodels.Promotion, error) {
	var out []promotionmodels.Promotion
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) CreatePromotion(ctx context.Context, p *promotionmodels.Promotion) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) UpdatePromotion(ctx context.Context, id string, fields map[string]any) error {
	res := r.db.WithContext(ctx).Model(&promotionmodels.Promotion{}).Where("id = ?", id).Updates(fields)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) DeletePromotion(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&promotionmodels.Promotion{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) ListCoupons(ctx context.Context) ([]promotionmodels.Coupon, error) {
	var out []promotionmodels.Coupon
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) CreateCoupon(ctx context.Context, c *promotionmodels.Coupon) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repo) UpdateCoupon(ctx context.Context, id string, fields map[string]any) error {
	res := r.db.WithContext(ctx).Model(&promotionmodels.Coupon{}).Where("id = ?", id).Updates(fields)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) DeleteCoupon(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&promotionmodels.Coupon{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
