package config

import (
	"time"
)

type App struct {
	Addr        string
	DataDir     string
	HTTPTimeout time.Duration

	Svc Service
	Dl  Downloader
}
