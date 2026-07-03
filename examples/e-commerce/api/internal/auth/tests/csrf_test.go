package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/auth"
	"github.com/yasinatesim/vela-commerce/api/internal/constants"
)

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(auth.RequireCSRF())
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestRequireCSRF_GetNeedsNoToken(t *testing.T) {
	r := newRouter()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireCSRF_PostWithoutTokenRejected(t *testing.T) {
	r := newRouter()
	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequireCSRF_PostWithMatchingTokenAccepted(t *testing.T) {
	r := newRouter()
	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	req.AddCookie(&http.Cookie{Name: constants.CSRF_COOKIE, Value: "tok123"})
	req.Header.Set(constants.CSRF_HEADER, "tok123")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireCSRF_PostWithMismatchedTokenRejected(t *testing.T) {
	r := newRouter()
	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	req.AddCookie(&http.Cookie{Name: constants.CSRF_COOKIE, Value: "tok123"})
	req.Header.Set(constants.CSRF_HEADER, "different")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
