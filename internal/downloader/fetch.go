package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"test-task/internal/util"
	"time"
)

const (
	renameRetries = 3
	renameDelay   = 50 * time.Millisecond
)

func (d *HTTPDownloader) Fetch(ctx context.Context, rawURL, suggestedName string) (string, int64, error) {
	u, err := parseURL(rawURL)
	if err != nil {
		return "", 0, err
	}

	dir := d.cfg.DownloadDir
	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", 0, err
	}

	name := d.pickFileName(suggestedName, u)
	final, part := targetPaths(dir, name)

	if n, ok, err := alreadyDone(final); err != nil {
		return "", 0, err
	} else if ok {
		return filepath.Base(final), n, nil
	}

	offset := partOffset(part)

	req, err := d.buildRequest(ctx, u.String(), offset)
	if err != nil {
		return "", 0, err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if offset, err = validateResumeStatus(resp.StatusCode, offset); err != nil {
		return "", 0, err
	}
	if err := validateContentType(resp.Header.Get("Content-Type")); err != nil {
		return "", 0, err
	}

	if disp := filenameFromDisposition(resp.Header); disp != "" && disp != name {
		name = disp
		final, part = targetPaths(dir, name)
		if n, ok, err := alreadyDone(final); err != nil {
			return "", 0, err
		} else if ok {
			return filepath.Base(final), n, nil
		}
		offset = partOffset(part)
	}

	f, err := openPart(part, offset)
	if err != nil {
		return "", 0, err
	}

	written, err := copyBody(f, resp.Body, d.cfg.MaxFileBytes, offset)
	if err != nil {
		_ = f.Close()
		return "", 0, err
	}

	if err := finalizePartFile(f, part, final); err != nil {
		return "", 0, err
	}
	return filepath.Base(final), offset + written, nil
}

func parseURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid url: %q", raw)
	}
	return u, nil
}

func (d *HTTPDownloader) pickFileName(suggested string, u *url.URL) string {
	if s := util.SanitizeFileName(strings.TrimSpace(suggested)); s != "" && s != "/" && s != "." {
		return s
	}
	if base := util.SanitizeFileName(path.Base(u.Path)); base != "" && base != "/" && base != "." {
		return base
	}
	return time.Now().Format("20060102_150405")
}

func targetPaths(dir, name string) (string, string) {
	final := filepath.Join(dir, name)
	return final, final + ".part"
}

func alreadyDone(final string) (int64, bool, error) {
	st, err := os.Stat(final)
	switch {
	case err == nil && st.Size() > 0:
		return st.Size(), true, nil
	case err != nil && !errors.Is(err, os.ErrNotExist):
		return 0, false, err
	default:
		return 0, false, nil
	}
}

func partOffset(part string) int64 {
	if st, err := os.Stat(part); err == nil {
		return st.Size()
	}
	return 0
}

func openPart(part string, offset int64) (*os.File, error) {
	if offset > 0 {
		return os.OpenFile(part, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	}
	return os.Create(part)
}

func (d *HTTPDownloader) buildRequest(ctx context.Context, urlStr string, offset int64) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	if offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	return req, nil
}

func validateResumeStatus(status int, offset int64) (int64, error) {
	switch status {
	case http.StatusOK:
		return 0, nil
	case http.StatusPartialContent:
		return offset, nil
	default:
		return 0, fmt.Errorf("unexpected status: %d", status)
	}
}

func validateContentType(ct string) error {
	if ct == "" {
		return nil
	}
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil
	}
	mt = strings.ToLower(mt)
	switch {
	case mt == "application/octet-stream",
		mt == "application/x-iso9660-image",
		mt == "application/zip",
		mt == "application/x-tar",
		strings.HasPrefix(mt, "image/"),
		strings.HasPrefix(mt, "video/"),
		strings.HasPrefix(mt, "audio/"):
		return nil
	case mt == "text/html", mt == "application/json":
		return fmt.Errorf("unexpected content-type: %s", mt)
	default:
		return nil
	}
}

func copyBody(dst io.Writer, src io.Reader, max, offset int64) (int64, error) {
	r := src
	if max > 0 {
		if remain := max - offset; remain <= 0 {
			return 0, fmt.Errorf("max file size reached")
		} else {
			r = io.LimitReader(src, remain)
		}
	}
	return io.Copy(dst, r)
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

func finalizePartFile(f *os.File, part, final string) error {
	_ = f.Sync()
	if err := f.Close(); err != nil {
		return err
	}
	_ = os.Remove(final)
	var err error
	for i := 0; i < renameRetries; i++ {
		if err = os.Rename(part, final); err == nil {
			return nil
		}
		time.Sleep(renameDelay)
	}
	return err
}
