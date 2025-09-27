package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken membuat random token 32 byte (64 hex char)
func GenerateSecureToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
