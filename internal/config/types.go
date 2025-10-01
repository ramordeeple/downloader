package config

import "time"

type Downloader struct {
	ClientTimeout time.Duration
	MaxFileBytes  int64
	MaxRetries    int
}

type Service struct {
	QueueSize int
	Workers   int
}
