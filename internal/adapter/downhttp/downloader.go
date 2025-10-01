package downhttp

import (
	"net/http"
	"test-task/internal/config"
)

type HTTPDownloader struct {
	client *http.Client
	cfg    config.Downloader
}

func New(cfg config.Downloader) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{Timeout: cfg.ClientTimeout},
		cfg:    cfg,
	}
}
