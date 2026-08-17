package returnreq

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	returnreqmodels "github.com/yasinatesim/vela-commerce/api/internal/returnreq/models"
)

type Handler struct {
	svc  *Service
	repo *Repository
}

func NewHandler(svc *Service, repo *Repository) *Handler {
	return &Handler{svc: svc, repo: repo}
}

type requestBody struct {
	OrderID string `json:"orderId" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
	Comment string `json:"comment"`
}

var statusByError = map[error]int{
	returnreqmodels.ErrOrderNotEligible:  http.StatusConflict,
	returnreqmodels.ErrReturnExists:      http.StatusConflict,
	returnreqmodels.ErrInvalidReason:     http.StatusBadRequest,
	returnreqmodels.ErrInvalidTransition: http.StatusConflict,
	returnreqmodels.ErrReturnNotFound:    http.StatusNotFound,
}

func httpStatus(err error) int {
	for known, code := range statusByError {
		if errors.Is(err, known) {
			return code
		}
	}
	return http.StatusInternalServerError
}

func (h *Handler) Create(c *gin.Context) {
	var body requestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	created, err := h.svc.Request(c.Request.Context(), body.OrderID, body.Reason, body.Comment)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *Handler) AdminList(c *gin.Context) {
	items, err := h.repo.List(c.Request.Context(), c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"returns": items})
}

type updateBody struct {
	Status    string `json:"status" binding:"required"`
	AdminNote string `json:"adminNote"`
}

func (h *Handler) AdminUpdateStatus(c *gin.Context) {
	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), c.Param("id"), body.Status, body.AdminNote); err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
