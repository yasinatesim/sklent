package order

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yasinatesim/vela-commerce/api/internal/auth"
	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/order/models"
	"github.com/yasinatesim/vela-commerce/api/internal/token"
)

var ErrInvalidStatus = errors.New("order: invalid status")

var ErrNotFound = errors.New("order: not found")

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) Create(ctx context.Context, o *ordermodels.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *Repo) GetTotalCents(ctx context.Context, orderID string) (int64, error) {
	var o ordermodels.Order
	err := r.db.WithContext(ctx).Select("total_cents").Where("id = ?", orderID).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrNotFound
	}
	return o.TotalCents, err
}

func (r *Repo) MarkPaid(ctx context.Context, orderID, paymentID string) error {
	return r.db.WithContext(ctx).Model(&ordermodels.Order{}).Where("id = ?", orderID).
		Updates(map[string]any{"status": constants.ORDER_STATUS_PAID, "payment_id": paymentID}).Error
}

func (r *Repo) GetByGuestToken(ctx context.Context, tok string) (*ordermodels.Order, error) {
	var o ordermodels.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("guest_token = ?", token.Hash(tok)).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &o, err
}

func (r *Repo) GetForUser(ctx context.Context, userID, orderID string) (*ordermodels.Order, error) {
	var o ordermodels.Order
	err := r.db.WithContext(ctx).Preload("Items").
		Where("id = ? AND user_id = ?", orderID, userID).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &o, err
}

func (r *Repo) ListForUser(ctx context.Context, userID string) ([]ordermodels.Order, error) {
	var out []ordermodels.Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) AdminList(ctx context.Context, status string, limit, offset int) ([]ordermodels.Order, error) {
	var out []ordermodels.Order
	q := r.db.WithContext(ctx)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Preload("Items").Order("created_at desc").Limit(limit).Offset(offset).Find(&out).Error
	return out, err
}

func (r *Repo) AdminUpdateStatus(ctx context.Context, orderID, status, trackingNumber string) error {
	validStatuses := map[string]bool{
		constants.ORDER_STATUS_PENDING:   true,
		constants.ORDER_STATUS_PAID:      true,
		constants.ORDER_STATUS_SHIPPED:   true,
		constants.ORDER_STATUS_CANCELLED: true,
	}
	if !validStatuses[status] {
		return ErrInvalidStatus
	}
	updates := map[string]any{"status": status}
	if trackingNumber != "" {
		updates["tracking_number"] = trackingNumber
	}
	return r.db.WithContext(ctx).Model(&ordermodels.Order{}).Where("id = ?", orderID).Updates(updates).Error
}

type Reserver interface {
	Reserve(ctx context.Context, orderID, productID string, qty int) error
}

type Mailer interface {
	SendOrderConfirmationAsync(orderID, email string, totalCents int64)
}

type Handler struct {
	repo     *Repo
	reserver Reserver
	mailer   Mailer
}

func NewHandler(repo *Repo, reserver Reserver, mailer Mailer) *Handler {
	return &Handler{repo: repo, reserver: reserver, mailer: mailer}
}

type placeItem struct {
	ProductID string `json:"productId" binding:"required"`
	TitleTr   string `json:"titleTr" binding:"max=200"`
	UnitCents int64  `json:"unitCents" binding:"required,min=1"`
	Quantity  int    `json:"quantity" binding:"required,min=1,max=99"`
}

type placeInput struct {
	Email         string      `json:"email" binding:"required,email"`
	PaymentMethod string      `json:"paymentMethod" binding:"required,oneof=card bank_transfer"`
	Items         []placeItem `json:"items" binding:"required,min=1,dive"`
}

// Place creates an order; members bind via OptionalAuth, guests get an unguessable token (hashed).
func (h *Handler) Place(c *gin.Context) {
	var in placeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}

	var total int64
	items := make([]ordermodels.OrderItem, 0, len(in.Items))
	for _, it := range in.Items {
		total += it.UnitCents * int64(it.Quantity)
		items = append(items, ordermodels.OrderItem{
			ProductID: it.ProductID, TitleTr: it.TitleTr, UnitCents: it.UnitCents, Quantity: it.Quantity,
		})
	}

	o := &ordermodels.Order{
		Email: in.Email, Status: constants.ORDER_STATUS_PENDING,
		PaymentMethod: in.PaymentMethod, TotalCents: total, Items: items,
	}

	rawToken := ""
	if userID, ok := auth.UserID(c); ok {
		o.UserID = &userID
	} else {
		raw, hash, err := token.Generate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		rawToken = raw
		o.GuestToken = hash
	}

	if err := h.repo.Create(c.Request.Context(), o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_failed"})
		return
	}

	for _, it := range o.Items {
		_ = h.reserver.Reserve(c.Request.Context(), o.ID, it.ProductID, it.Quantity)
	}
	h.mailer.SendOrderConfirmationAsync(o.ID, o.Email, o.TotalCents)

	c.JSON(http.StatusCreated, gin.H{"orderId": o.ID, "totalCents": total, "trackToken": rawToken})
}

func (h *Handler) GetByGuestToken(c *gin.Context) {
	o, err := h.repo.GetByGuestToken(c.Request.Context(), c.Param("token"))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListForUser(c *gin.Context) {
	userID, _ := auth.UserID(c)
	items, err := h.repo.ListForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, _ := auth.UserID(c)
	o, err := h.repo.GetForUser(c.Request.Context(), userID, c.Param("id"))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) AdminList(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit <= 0 || limit > 200 {
		limit = 50
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	items, err := h.repo.AdminList(c.Request.Context(), c.Query("status"), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type adminUpdateStatusInput struct {
	Status         string `json:"status" binding:"required,oneof=pending paid shipped cancelled"`
	TrackingNumber string `json:"trackingNumber" binding:"max=100"`
}

func (h *Handler) AdminUpdateStatus(c *gin.Context) {
	var in adminUpdateStatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	if err := h.repo.AdminUpdateStatus(c.Request.Context(), c.Param("id"), in.Status, in.TrackingNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// StatusAndEmail exposes just enough of an order for other domains (returns, invoicing) to make
// their own eligibility decisions without importing the order handler.
func (r *Repo) StatusAndEmail(ctx context.Context, orderID string) (string, string, error) {
	var row struct {
		Status string
		Email  string
	}
	err := r.db.WithContext(ctx).Model(&ordermodels.Order{}).
		Select("status", "email").
		Where("id = ?", orderID).
		Take(&row).Error
	if err != nil {
		return "", "", err
	}
	return row.Status, row.Email, nil
}
