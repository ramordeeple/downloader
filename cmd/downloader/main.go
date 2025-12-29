package main

import (
	"context"
	"downloader/internal/app"
	"downloader/internal/config"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.App{
		Addr:        ":8080",
		DataDir:     "data",
		HTTPTimeout: time.Second * 60,
		Svc: config.Service{
			QueueSize: 100,
			Workers:   2,
		},
		Dl: config.Downloader{
			ClientTimeout: 60 * time.Second,
			MaxRetries:    2,
		},
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	a := app.New(cfg)
	a.Start()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = a.Shutdown(shutdownCtx)
}
