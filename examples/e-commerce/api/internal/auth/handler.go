package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/token"
	"github.com/yasinatesim/vela-commerce/api/internal/user/models"
)

const refreshTTL = 7 * 24 * time.Hour
const refreshCookie = "refresh_token"

type Handler struct {
	signer       *Signer
	repo         *Repo
	cookieDomain string
	secure       bool
}

func NewHandler(signer *Signer, repo *Repo, cookieDomain string, secure bool) *Handler {
	return &Handler{signer: signer, repo: repo, cookieDomain: cookieDomain, secure: secure}
}

type credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	FullName string `json:"fullName" binding:"max=120"`
}

func (h *Handler) Register(c *gin.Context) {
	var in credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	hash, err := HashPassword(in.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	u := &usermodels.User{Email: in.Email, PasswordHash: hash, Role: constants.ROLE_USER, FullName: in.FullName}
	if err := h.repo.CreateUser(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email_taken"})
		return
	}
	h.issueSession(c, u)
}

func (h *Handler) Login(c *gin.Context) {
	var in credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
		return
	}
	u, err := h.repo.FindByEmail(c.Request.Context(), in.Email)
	if err != nil || !VerifyPassword(u.PasswordHash, in.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
		return
	}
	if u.ClosedAt != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "account_closed"})
		return
	}
	h.issueSession(c, u)
}

func (h *Handler) Me(c *gin.Context) {
	claims, err := h.signer.Parse(h.signer.bearerOrCookie(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userId": claims.UserID, "role": claims.Role})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(accessCookie, "", -1, "/", h.cookieDomain, h.secure, true)
	c.SetCookie(refreshCookie, "", -1, "/", h.cookieDomain, h.secure, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) Refresh(c *gin.Context) {
	raw, err := c.Cookie(refreshCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no_refresh"})
		return
	}
	userID, err := h.repo.RotateRefresh(c.Request.Context(), raw)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh"})
		return
	}
	u, err := h.repo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh"})
		return
	}
	h.issueSession(c, u)
}

func (h *Handler) issueSession(c *gin.Context, u *usermodels.User) {
	access, err := h.signer.IssueAccess(u.ID, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	rawRefresh, _, err := token.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if err := h.repo.StoreRefresh(c.Request.Context(), u.ID, rawRefresh, refreshTTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.SetCookie(accessCookie, access, int(accessTTL.Seconds()), "/", h.cookieDomain, h.secure, true)
	c.SetCookie(refreshCookie, rawRefresh, int(refreshTTL.Seconds()), "/", h.cookieDomain, h.secure, true)
	c.JSON(http.StatusOK, gin.H{"userId": u.ID, "role": u.Role, "email": u.Email})
}
