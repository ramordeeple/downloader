package downloader

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"test-task/internal/util"
)

func (d *HTTPDownloader) tryFetch(ctx context.Context, url, name string) (string, int64, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", 0, errors.New(resp.Status)
	}

	filename := util.UniquePath(filepath.Join(d.cfg.DownloadDir, name))
	f, err := os.Create(filename)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	r := io.Reader(resp.Body)
	if d.cfg.MaxFileBytes > 0 {
		r = io.LimitReader(resp.Body, d.cfg.MaxFileBytes)
	}
	n, err := io.Copy(f, r)
	if err != nil {
		_ = os.Remove(filename)
		return "", 0, err
	}

	return filepath.Base(filename), n, nil
}

func (d *HTTPDownloader) Fetch(ctx context.Context, rawURL, suggestedName string) (string, int64, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return "", 0, err
	}

	name := fileName(u, suggestedName)
	if err := os.MkdirAll(d.cfg.DownloadDir, 0o755); err != nil {
		return "", 0, err
	}

	var lastErr error
	for attempt := 0; attempt < d.cfg.MaxRetries; attempt++ {
		if attempt > 0 && !sleepOrCancel(ctx, attempt) {
			return "", 0, ctx.Err()
		}
		if saved, n, err := d.tryFetch(ctx, u.String(), name); err == nil {
			return saved, n, nil
		} else {
			lastErr = err
		}
	}

	return "", 0, lastErr
}
