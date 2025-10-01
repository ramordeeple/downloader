package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandID(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
