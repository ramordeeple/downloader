package main

import (
	"context"
	"os/signal"
	"syscall"
	"test-task/internal/app"
	"test-task/internal/service"
	"time"
)

func main() {
	cfg := app.Config{
		Addr:        ":8080",
		DataDir:     "data",
		DownloadDir: "downloads",
		HTTPTimeout: time.Second * 60,
		MaxRetries:  2,
		Svc: service.Config{
			QueueSize: 200,
			Workers:   4,
		},
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	a := app.New(cfg)
	a.Start(ctx)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = a.Shutdown(shutdownCtx)
}
