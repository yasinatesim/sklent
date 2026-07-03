package promotion

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	promotionmodels "github.com/yasinatesim/vela-commerce/api/internal/promotion/models"
)

func joinIDs(ids []string) string { return strings.Join(ids, ",") }

func splitIDs(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler { return &Handler{repo: repo} }

type promotionInput struct {
	Name          string   `json:"name" binding:"required,max=120"`
	DiscountType  string   `json:"discountType" binding:"required,oneof=percent fixed_try"`
	DiscountValue int      `json:"discountValue" binding:"required,min=1"`
	ScopeType     string   `json:"scopeType" binding:"required,oneof=all products categories"`
	ProductIDs    []string `json:"productIds"`
	CategoryIDs   []string `json:"categoryIds"`
	MinCartCents  int64    `json:"minCartCents" binding:"min=0"`
}

func promotionToJSON(p promotionmodels.Promotion) gin.H {
	return gin.H{
		"id": p.ID, "name": p.Name, "discountType": p.DiscountType, "discountValue": p.DiscountValue,
		"scopeType": p.ScopeType, "productIds": splitIDs(p.ProductIDs), "categoryIds": splitIDs(p.CategoryIDs),
		"minCartCents": p.MinCartCents, "active": p.Active,
	}
}

func (h *Handler) AdminListPromotions(c *gin.Context) {
	items, err := h.repo.ListPromotions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := make([]gin.H, 0, len(items))
	for _, p := range items {
		out = append(out, promotionToJSON(p))
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

func (h *Handler) AdminCreatePromotion(c *gin.Context) {
	var in promotionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	p := &promotionmodels.Promotion{
		Name: in.Name, DiscountType: in.DiscountType, DiscountValue: in.DiscountValue,
		ScopeType: in.ScopeType, ProductIDs: joinIDs(in.ProductIDs), CategoryIDs: joinIDs(in.CategoryIDs),
		MinCartCents: in.MinCartCents, Active: true,
	}
	if err := h.repo.CreatePromotion(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_failed"})
		return
	}
	c.JSON(http.StatusCreated, promotionToJSON(*p))
}

type promotionPatchInput struct {
	Name          *string  `json:"name"`
	DiscountType  *string  `json:"discountType" binding:"omitempty,oneof=percent fixed_try"`
	DiscountValue *int     `json:"discountValue"`
	ScopeType     *string  `json:"scopeType" binding:"omitempty,oneof=all products categories"`
	ProductIDs    []string `json:"productIds"`
	CategoryIDs   []string `json:"categoryIds"`
	MinCartCents  *int64   `json:"minCartCents"`
	Active        *bool    `json:"active"`
}

func (h *Handler) AdminUpdatePromotion(c *gin.Context) {
	var in promotionPatchInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	fields := map[string]any{}
	if in.Name != nil {
		fields["name"] = *in.Name
	}
	if in.DiscountType != nil {
		fields["discount_type"] = *in.DiscountType
	}
	if in.DiscountValue != nil {
		fields["discount_value"] = *in.DiscountValue
	}
	if in.ScopeType != nil {
		fields["scope_type"] = *in.ScopeType
	}
	if in.ProductIDs != nil {
		fields["product_ids"] = joinIDs(in.ProductIDs)
	}
	if in.CategoryIDs != nil {
		fields["category_ids"] = joinIDs(in.CategoryIDs)
	}
	if in.MinCartCents != nil {
		fields["min_cart_cents"] = *in.MinCartCents
	}
	if in.Active != nil {
		fields["active"] = *in.Active
	}
	if err := h.repo.UpdatePromotion(c.Request.Context(), c.Param("id"), fields); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminDeletePromotion(c *gin.Context) {
	if err := h.repo.DeletePromotion(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type couponInput struct {
	Code          string   `json:"code" binding:"required,max=40"`
	DiscountType  string   `json:"discountType" binding:"required,oneof=percent fixed_try"`
	DiscountValue int      `json:"discountValue" binding:"required,min=1"`
	ScopeType     string   `json:"scopeType" binding:"required,oneof=all products categories"`
	ProductIDs    []string `json:"productIds"`
	CategoryIDs   []string `json:"categoryIds"`
	MinCartCents  int64    `json:"minCartCents" binding:"min=0"`
}

func couponToJSON(c promotionmodels.Coupon) gin.H {
	return gin.H{
		"id": c.ID, "code": c.Code, "discountType": c.DiscountType, "discountValue": c.DiscountValue,
		"scopeType": c.ScopeType, "productIds": splitIDs(c.ProductIDs), "categoryIds": splitIDs(c.CategoryIDs),
		"minCartCents": c.MinCartCents, "active": c.Active,
	}
}

func (h *Handler) AdminListCoupons(c *gin.Context) {
	items, err := h.repo.ListCoupons(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := make([]gin.H, 0, len(items))
	for _, cp := range items {
		out = append(out, couponToJSON(cp))
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

func (h *Handler) AdminCreateCoupon(c *gin.Context) {
	var in couponInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	cp := &promotionmodels.Coupon{
		Code: strings.ToUpper(strings.TrimSpace(in.Code)), DiscountType: in.DiscountType, DiscountValue: in.DiscountValue,
		ScopeType: in.ScopeType, ProductIDs: joinIDs(in.ProductIDs), CategoryIDs: joinIDs(in.CategoryIDs),
		MinCartCents: in.MinCartCents, Active: true,
	}
	if err := h.repo.CreateCoupon(c.Request.Context(), cp); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "code_taken"})
		return
	}
	c.JSON(http.StatusCreated, couponToJSON(*cp))
}

type couponPatchInput struct {
	DiscountType  *string  `json:"discountType" binding:"omitempty,oneof=percent fixed_try"`
	DiscountValue *int     `json:"discountValue"`
	ScopeType     *string  `json:"scopeType" binding:"omitempty,oneof=all products categories"`
	ProductIDs    []string `json:"productIds"`
	CategoryIDs   []string `json:"categoryIds"`
	MinCartCents  *int64   `json:"minCartCents"`
	Active        *bool    `json:"active"`
}

func (h *Handler) AdminUpdateCoupon(c *gin.Context) {
	var in couponPatchInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	fields := map[string]any{}
	if in.DiscountType != nil {
		fields["discount_type"] = *in.DiscountType
	}
	if in.DiscountValue != nil {
		fields["discount_value"] = *in.DiscountValue
	}
	if in.ScopeType != nil {
		fields["scope_type"] = *in.ScopeType
	}
	if in.ProductIDs != nil {
		fields["product_ids"] = joinIDs(in.ProductIDs)
	}
	if in.CategoryIDs != nil {
		fields["category_ids"] = joinIDs(in.CategoryIDs)
	}
	if in.MinCartCents != nil {
		fields["min_cart_cents"] = *in.MinCartCents
	}
	if in.Active != nil {
		fields["active"] = *in.Active
	}
	if err := h.repo.UpdateCoupon(c.Request.Context(), c.Param("id"), fields); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminDeleteCoupon(c *gin.Context) {
	if err := h.repo.DeleteCoupon(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
