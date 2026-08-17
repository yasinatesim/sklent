package constants_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yasinatesim/vela-commerce/api/internal/constants"
)

func TestFrontendPathsAreLocaleNeutral(t *testing.T) {
	paths := map[string]string{
		"checkout success": constants.FRONTEND_PATH_CHECKOUT_SUCCESS,
		"checkout error":   constants.FRONTEND_PATH_CHECKOUT_ERROR,
	}

	turkish := []string{"odeme", "basarili", "hata", "sepet", "siparis", "urun", "kategori"}

	for name, path := range paths {
		t.Run(name, func(t *testing.T) {
			assert.True(t, strings.HasPrefix(path, "/"), "must be an absolute path")
			assert.False(t, strings.HasPrefix(path, "//"), "must not be protocol-relative")
			for _, segment := range strings.Split(strings.Trim(path, "/"), "/") {
				assert.NotContains(t, turkish, segment, "route segments stay locale-neutral")
				assert.Equal(t, strings.ToLower(segment), segment, "segments are lowercase")
			}
		})
	}
}
