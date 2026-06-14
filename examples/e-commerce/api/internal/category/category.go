package category

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/category/models"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context) ([]categorymodels.Category, error) {
	var out []categorymodels.Category
	err := r.db.WithContext(ctx).Order("name_tr asc").Find(&out).Error
	return out, err
}

func (r *Repo) Create(ctx context.Context, c *categorymodels.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

type Handler struct{ repo *Repo }

func NewHandler(repo *Repo) *Handler { return &Handler{repo: repo} }

func (h *Handler) List(c *gin.Context) {
	items, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
