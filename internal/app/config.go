package app

import (
	"test-task/internal/service"
	"time"
)

type Config struct {
	Addr        string
	DataDir     string
	DownloadDir string
	HTTPTimeout time.Duration
	MaxRetries  int
	Svc         service.Config
}
