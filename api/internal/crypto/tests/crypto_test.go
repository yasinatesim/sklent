package crypto_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasinatesim/vela-commerce/api/internal/crypto"
)

func key32() string {
	return base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	c, err := crypto.New(key32())
	require.NoError(t, err)

	enc, err := c.Encrypt("sk-secret-provider-key")
	require.NoError(t, err)
	assert.NotEqual(t, "sk-secret-provider-key", enc)

	dec, err := c.Decrypt(enc)
	require.NoError(t, err)
	assert.Equal(t, "sk-secret-provider-key", dec)
}

func TestEncryptIsNonDeterministic(t *testing.T) {
	c, err := crypto.New(key32())
	require.NoError(t, err)
	a, _ := c.Encrypt("same")
	b, _ := c.Encrypt("same")
	assert.NotEqual(t, a, b)
}

func TestRejectsWrongKeySize(t *testing.T) {
	short := base64.StdEncoding.EncodeToString([]byte("too-short"))
	_, err := crypto.New(short)
	assert.ErrorIs(t, err, crypto.ErrKeySize)
}
