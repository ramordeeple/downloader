package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SanitizeFileName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "..", "_")

	return s
}

func UniquePath(p string) string {
	base := p
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return p
	}
	for i := 1; i < 10_000; i++ {
		cand := fmt.Sprintf("%s(%d)%s", strings.TrimSuffix(base, filepath.Ext(base)), i, filepath.Ext(base))
		if _, err := os.Stat(cand); errors.Is(err, os.ErrNotExist) {
			return cand
		}
	}
	return base + ".dup"
}
