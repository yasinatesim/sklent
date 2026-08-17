package product

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/product/models"
	"github.com/yasinatesim/vela-commerce/api/internal/rag"
)

var ErrNotFound = errors.New("product: not found")

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListPublished(ctx context.Context, categorySlug string) ([]productmodels.Product, error) {
	var out []productmodels.Product
	q := r.db.WithContext(ctx).Where("published = ?", true)
	if categorySlug != "" && categorySlug != "all" {
		q = q.Where("category_slug = ?", categorySlug)
	}
	err := q.Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) GetBySlug(ctx context.Context, slug string) (*productmodels.Product, error) {
	var p productmodels.Product
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *Repo) Create(ctx context.Context, p *productmodels.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) GetIDBySlug(ctx context.Context, slug string) (string, error) {
	p, err := r.GetBySlug(ctx, slug)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

// DecrementStock is atomic: it only decrements if enough stock remains, then returns the resulting state for alerting.
func (r *Repo) DecrementStock(ctx context.Context, productID string, qty int) (newStock, threshold int, title string, err error) {
	res := r.db.WithContext(ctx).Model(&productmodels.Product{}).
		Where("id = ? AND stock >= ?", productID, qty).
		Update("stock", gorm.Expr("stock - ?", qty))
	if res.Error != nil {
		return 0, 0, "", res.Error
	}
	var p productmodels.Product
	if err := r.db.WithContext(ctx).Select("stock", "low_stock_threshold", "title_tr").
		Where("id = ?", productID).First(&p).Error; err != nil {
		return 0, 0, "", err
	}
	return p.Stock, p.LowStockThreshold, p.TitleTr, nil
}

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler { return &Handler{repo: repo} }

func (h *Handler) List(c *gin.Context) {
	q := ParseListQuery(c.Request.URL.Query())
	items, pagination, err := h.repo.List(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *Handler) Facets(c *gin.Context) {
	facets, err := h.repo.Facets(c.Request.Context(), ParseListQuery(c.Request.URL.Query()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, facets)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	p, err := h.repo.GetBySlug(c.Request.Context(), c.Param("slug"))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, p)
}

type createInput struct {
	Slug          string `json:"slug" binding:"max=160"`
	TitleTr       string `json:"titleTr" binding:"required,max=200"`
	TitleEn       string `json:"titleEn" binding:"max=200"`
	DescriptionTr string `json:"descriptionTr" binding:"max=2000"`
	PriceCents    int64  `json:"priceCents" binding:"required,min=1"`
	OldPriceCents int64  `json:"oldPriceCents" binding:"min=0"`
	Stock         int    `json:"stock" binding:"min=0"`
	CategorySlug  string `json:"categorySlug" binding:"max=64"`
	Badge         string `json:"badge" binding:"max=40"`
	Seller        string `json:"seller" binding:"max=80"`
	Published     bool   `json:"published"`
	ImageURL      string `json:"imageUrl" binding:"max=500"`
}

func (h *Handler) AdminCreate(c *gin.Context) {
	var in createInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	slug := in.Slug
	if slug == "" {
		slug = rag.Slugify(in.TitleTr)
	}
	p := &productmodels.Product{
		Slug: slug, TitleTr: in.TitleTr, TitleEn: in.TitleEn, DescriptionTr: in.DescriptionTr,
		PriceCents: in.PriceCents, OldPriceCents: in.OldPriceCents, Stock: in.Stock,
		CategorySlug: in.CategorySlug, Badge: in.Badge, Seller: in.Seller,
		Published: in.Published, ImageURL: in.ImageURL,
	}
	if err := h.repo.Create(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_failed"})
		return
	}
	c.JSON(http.StatusCreated, p)
}
