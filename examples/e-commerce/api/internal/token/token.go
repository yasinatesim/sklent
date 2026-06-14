package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Generate returns a random opaque token (raw, shown once) and its sha256 hash (stored).
func Generate() (raw string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(buf)
	hash = Hash(raw)
	return raw, hash, nil
}

func Hash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
