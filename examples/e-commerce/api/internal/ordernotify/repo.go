// Package ordernotify mails the admin-configured recipients when an order arrives from the site or a marketplace.
package ordernotify

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/ordernotify/models"
)

var ErrNotFound = errors.New("order notification rule not found")

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context) ([]ordernotifymodels.Rule, error) {
	rules := []ordernotifymodels.Rule{}
	err := r.db.WithContext(ctx).Order("source ASC, recipient ASC").Find(&rules).Error
	return rules, err
}

// ListEnabledBySource returns the rules that should fire for one channel.
func (r *Repo) ListEnabledBySource(ctx context.Context, source string) ([]ordernotifymodels.Rule, error) {
	rules := []ordernotifymodels.Rule{}
	err := r.db.WithContext(ctx).
		Where("enabled AND source = ?", source).
		Order("recipient ASC").Find(&rules).Error
	return rules, err
}

func (r *Repo) Create(ctx context.Context, rule *ordernotifymodels.Rule) error {
	rule.Recipient = strings.TrimSpace(rule.Recipient)
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *Repo) Update(ctx context.Context, id string, fields map[string]any) error {
	res := r.db.WithContext(ctx).Model(&ordernotifymodels.Rule{}).Where("id = ?", id).Updates(fields)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&ordernotifymodels.Rule{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
