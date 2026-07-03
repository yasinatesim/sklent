package review

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/review/models"
)

var ErrNotFound = errors.New("review: not found")

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListApprovedForProduct(ctx context.Context, productID string) ([]reviewmodels.Review, error) {
	var out []reviewmodels.Review
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND status = ?", productID, constants.REVIEW_STATUS_APPROVED).
		Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) Create(ctx context.Context, productID, authorName, comment string, rating int) (*reviewmodels.Review, error) {
	rv := &reviewmodels.Review{
		ProductID: productID, AuthorName: authorName, Rating: rating, Comment: comment,
		Status: constants.REVIEW_STATUS_PENDING,
	}
	if err := r.db.WithContext(ctx).Create(rv).Error; err != nil {
		return nil, err
	}
	return rv, nil
}

func (r *Repo) AdminList(ctx context.Context) ([]reviewmodels.Review, error) {
	var out []reviewmodels.Review
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) AdminUpdateStatus(ctx context.Context, id, status string) error {
	res := r.db.WithContext(ctx).Model(&reviewmodels.Review{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// SlugResolver looks up a product's ID from its public slug without review depending on the product package directly.
type SlugResolver interface {
	GetIDBySlug(ctx context.Context, slug string) (string, error)
}

type Handler struct {
	repo     *Repo
	resolver SlugResolver
}

func NewHandler(repo *Repo, resolver SlugResolver) *Handler {
	return &Handler{repo: repo, resolver: resolver}
}

func (h *Handler) List(c *gin.Context) {
	productID, err := h.resolver.GetIDBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	items, err := h.repo.ListApprovedForProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createInput struct {
	AuthorName string `json:"authorName" binding:"required,max=80"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Comment    string `json:"comment" binding:"max=2000"`
}

func (h *Handler) Create(c *gin.Context) {
	productID, err := h.resolver.GetIDBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in createInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	rv, err := h.repo.Create(c.Request.Context(), productID, in.AuthorName, in.Comment, in.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_failed"})
		return
	}
	c.JSON(http.StatusCreated, rv)
}

func (h *Handler) AdminList(c *gin.Context) {
	items, err := h.repo.AdminList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type updateStatusInput struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

func (h *Handler) AdminUpdateStatus(c *gin.Context) {
	var in updateStatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	if err := h.repo.AdminUpdateStatus(c.Request.Context(), c.Param("id"), in.Status); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
