package downhttp

import (
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"test-task/internal/util"
)

func (d *HTTPDownloader) pickFileName(suggested string, u *url.URL) string {
	if s := util.SanitizeFileName(strings.TrimSpace(suggested)); s != "" && s != "/" && s != "." {
		return s
	}
	base := util.SanitizeFileName(path.Base(u.Path))
	if base == "" || base == "/" || base == "." {
		base = time.Now().Format("20060102_150405")
	}
	return base
}

func filenameFromDisposition(h http.Header) string {
	cd := h.Get("Content-Disposition")
	if cd == "" {
		return ""
	}
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return ""
	}
	if v := params["filename*"]; v != "" {
		if i := strings.LastIndex(v, "''"); i >= 0 && i+2 < len(v) {
			return util.SanitizeFileName(v[i+2:])
		}
		return util.SanitizeFileName(v)
	}
	if v := params["filename"]; v != "" {
		return util.SanitizeFileName(v)
	}
	return ""
}
