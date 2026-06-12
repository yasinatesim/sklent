package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
	"github.com/yasinatesim/vela-commerce/api/internal/token"
)

// IssueCSRF sets a readable cookie the client echoes back in X-CSRF-Token (double-submit).
func IssueCSRF(cookieDomain string, secure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Cookie(constants.CSRF_COOKIE); err != nil {
			raw, _, _ := token.Generate()
			c.SetCookie(constants.CSRF_COOKIE, raw, 3600, "/", cookieDomain, secure, false)
		}
		c.Next()
	}
}

func RequireCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(constants.CSRF_COOKIE)
		header := c.GetHeader(constants.CSRF_HEADER)
		if err != nil || cookie == "" || header == "" || cookie != header {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "csrf_mismatch"})
			return
		}
		c.Next()
	}
}
