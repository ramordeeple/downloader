package downloader

import (
	"net/url"
	"path"
	"strings"
	"test-task/internal/util"
	"time"
)

func fileName(u *url.URL, suggestedName string) string {
	name := util.SanitizeFileName(strings.TrimSpace(suggestedName))

	if name == "" {
		name = path.Base(u.Path)
	}
	if name == "" || name == "/" || name == "." {
		name = time.Now().Format("20060102_150405")
	}

	return name
}
