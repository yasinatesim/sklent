package ordernotify

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yasinatesim/vela-commerce/api/internal/ordernotify/models"
)

type Handler struct {
	repo *Repo
	svc  *Service
}

func NewHandler(repo *Repo, svc *Service) *Handler { return &Handler{repo: repo, svc: svc} }

type ruleInput struct {
	Source    string `json:"source"`
	Recipient string `json:"recipient"`
	Enabled   *bool  `json:"enabled"`
}

// isPlausibleEmail keeps obvious junk out; real deliverability is proven by the send itself.
func isPlausibleEmail(v string) bool {
	at := strings.IndexByte(v, '@')
	return at > 0 && at < len(v)-1 && strings.Contains(v[at+1:], ".") && !strings.ContainsAny(v, " \t\r\n")
}

// AdminList handles GET /admin/order-notifications.
func (h *Handler) AdminList(c *gin.Context) {
	rules, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rules})
}

// AdminCreate handles POST /admin/order-notifications.
func (h *Handler) AdminCreate(c *gin.Context) {
	var in ruleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	recipient := strings.TrimSpace(in.Recipient)
	if !isPlausibleEmail(recipient) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid_recipient"})
		return
	}
	source := strings.TrimSpace(in.Source)
	if !IsKnownSource(source) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "unknown_source"})
		return
	}
	rule := ordernotifymodels.Rule{Source: source, Recipient: recipient, Enabled: true}
	if in.Enabled != nil {
		rule.Enabled = *in.Enabled
	}
	if err := h.repo.Create(c.Request.Context(), &rule); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "create_failed"})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

// AdminUpdate handles PATCH /admin/order-notifications/:id.
func (h *Handler) AdminUpdate(c *gin.Context) {
	var in ruleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	fields := map[string]any{}
	if recipient := strings.TrimSpace(in.Recipient); recipient != "" {
		if !isPlausibleEmail(recipient) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid_recipient"})
			return
		}
		fields["recipient"] = recipient
	}
	if source := strings.TrimSpace(in.Source); source != "" {
		if !IsKnownSource(source) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "unknown_source"})
			return
		}
		fields["source"] = source
	}
	if in.Enabled != nil {
		fields["enabled"] = *in.Enabled
	}
	if len(fields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no_fields"})
		return
	}
	if err := h.repo.Update(c.Request.Context(), c.Param("id"), fields); err != nil {
		h.writeRuleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// AdminDelete handles DELETE /admin/order-notifications/:id.
func (h *Handler) AdminDelete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.writeRuleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) writeRuleError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
}
