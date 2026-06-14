package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
)

const (
	ctxUserID = "auth.userID"
	ctxRole   = "auth.role"

	accessCookie = "access_token"
)

func (s *Signer) bearerOrCookie(c *gin.Context) string {
	if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	if cookie, err := c.Cookie(accessCookie); err == nil {
		return cookie
	}
	return ""
}

func (s *Signer) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := s.Parse(s.bearerOrCookie(c))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

func (s *Signer) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, err := s.Parse(s.bearerOrCookie(c)); err == nil {
			c.Set(ctxUserID, claims.UserID)
			c.Set(ctxRole, claims.Role)
		}
		c.Next()
	}
}

func (s *Signer) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := s.Parse(s.bearerOrCookie(c))
		if err != nil || claims.Role != constants.ROLE_ADMIN {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

func UserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return "", false
	}
	id, _ := v.(string)
	return id, id != ""
}
