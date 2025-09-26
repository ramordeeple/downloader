package downloader

import (
	"net/http"
	"time"
)

type HTTPDownloader struct {
	cfg    Config
	client *http.Client
}

func NewHTTPDownloader(cfg Config) *HTTPDownloader {
	if cfg.ClientTimeout == 0 {
		cfg.ClientTimeout = time.Second * 60
	}

	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}

	return &HTTPDownloader{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.ClientTimeout},
	}
}
