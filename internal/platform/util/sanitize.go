package util

import "strings"

// SanitizeFileName очищает имя файла от опасных символов.
func SanitizeFileName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "..", "_")
	s = strings.ReplaceAll(s, "\"", "_")
	return s
}
