package app

import (
	"context"
	"log"
	"net/http"
	"test-task/internal/api"
	"test-task/internal/downloader"
	"test-task/internal/service"
	"test-task/internal/store"
	"test-task/internal/util"
)

type App struct {
	server  *http.Server
	service *service.Service
}

func New(cfg Config) *App {
	st := store.NewFileSys(cfg.DataDir)
	dl := downloader.NewHTTPDownloader(downloader.Config{
		DownloadDir:   cfg.DownloadDir,
		ClientTimeout: cfg.HTTPTimeout,
	})

	s := service.New(st, dl, cfg.Svc)
	s.Load()

	mux := api.NewMux(api.NewAPI(s))
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: util.LogRequests(mux),
	}
	return &App{server: srv, service: s}
}

func (app *App) Start(ctx context.Context) {
	app.service.Start(ctx, 4)
	go func() {
		log.Printf("listening on %s", app.server.Addr)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
