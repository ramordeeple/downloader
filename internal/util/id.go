package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandID(n int) string {
	if n <= 0 {
		n = 6
	}
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
