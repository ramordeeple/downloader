package downloader

import "time"

type Config struct {
	DownloadDir   string
	ClientTimeout time.Duration
	MaxFileBytes  int64
	MaxRetries    int
}
