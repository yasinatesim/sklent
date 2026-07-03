package stocktracking

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	stocktrackingmodels "github.com/yasinatesim/vela-commerce/api/internal/stocktracking/models"
)

var ErrNotFound = errors.New("stocktracking: not found")

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context) ([]stocktrackingmodels.StockTrackingItem, error) {
	var out []stocktrackingmodels.StockTrackingItem
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) Create(ctx context.Context, productName string, quantity int) (*stocktrackingmodels.StockTrackingItem, error) {
	item := &stocktrackingmodels.StockTrackingItem{ProductName: productName, Quantity: quantity}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repo) Update(ctx context.Context, id, productName string, quantity int) error {
	res := r.db.WithContext(ctx).Model(&stocktrackingmodels.StockTrackingItem{}).Where("id = ?", id).
		Updates(map[string]any{"product_name": productName, "quantity": quantity})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&stocktrackingmodels.StockTrackingItem{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler { return &Handler{repo: repo} }

func (h *Handler) AdminList(c *gin.Context) {
	items, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createInput struct {
	ProductName string `json:"productName" binding:"required,max=200"`
	Quantity    int    `json:"quantity" binding:"min=0"`
}

func (h *Handler) AdminCreate(c *gin.Context) {
	var in createInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	item, err := h.repo.Create(c.Request.Context(), in.ProductName, in.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_failed"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) AdminUpdate(c *gin.Context) {
	var in createInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	if err := h.repo.Update(c.Request.Context(), c.Param("id"), in.ProductName, in.Quantity); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminDelete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
