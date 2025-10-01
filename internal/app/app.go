package app

import (
	"context"
	"net/http"
	"sync"
	"time"

	"test-task/internal/adapter/downhttp"
	"test-task/internal/adapter/httpapi"
	"test-task/internal/adapter/storefs"
	"test-task/internal/config"
	"test-task/internal/platform/queue"
	"test-task/internal/usecase"
)

type App struct {
	srv     *http.Server
	uc      *usecase.TaskService
	workers int
}

type sysClock struct{}

func (sysClock) Now() time.Time { return time.Now() }

func New(cfg config.App) *App {
	repo := storefs.New(cfg.DataDir)
	fetcher := downhttp.New(cfg.Dl)
	q := queue.New(cfg.Svc.QueueSize)
	clock := sysClock{}

	uc := usecase.NewTaskService(
		repo,
		fetcher,
		q,
		clock,
		cfg.DataDir,
	)

	mux := httpapi.New(uc)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &App{
		srv:     srv,
		uc:      uc,
		workers: cfg.Svc.Workers,
	}
}

func (a *App) Start() {
	_ = a.uc.Restore()

	a.uc.Start(a.workers)

	go func() { _ = a.srv.ListenAndServe() }()
}

func (a *App) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)
	a.uc.Stop(&wg)
	wg.Wait()
	return a.srv.Shutdown(ctx)
}
