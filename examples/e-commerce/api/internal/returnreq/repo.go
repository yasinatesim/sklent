package returnreq

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	returnreqmodels "github.com/yasinatesim/vela-commerce/api/internal/returnreq/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, ret returnreqmodels.Return) (returnreqmodels.Return, error) {
	if err := r.db.WithContext(ctx).Create(&ret).Error; err != nil {
		return returnreqmodels.Return{}, err
	}
	return ret, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*returnreqmodels.Return, error) {
	var out returnreqmodels.Return
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, returnreqmodels.ErrReturnNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// OpenForOrder returns the return still in flight for an order, so a second request cannot be
// opened while one is pending.
func (r *Repository) OpenForOrder(ctx context.Context, orderID string) (*returnreqmodels.Return, error) {
	var out returnreqmodels.Return
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND status IN ?", orderID,
			[]string{constants.RETURN_STATUS_REQUESTED, constants.RETURN_STATUS_APPROVED}).
		First(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) List(ctx context.Context, status string) ([]returnreqmodels.Return, error) {
	var out []returnreqmodels.Return
	q := r.db.WithContext(ctx).Order("created_at DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id, status, adminNote string) error {
	res := r.db.WithContext(ctx).Model(&returnreqmodels.Return{}).
		Where("id = ?", id).
		Updates(map[string]any{"status": status, "admin_note": adminNote})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return returnreqmodels.ErrReturnNotFound
	}
	return nil
}
