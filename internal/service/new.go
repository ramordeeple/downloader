package service

import (
	"context"
	"test-task/internal/downloader"
	"test-task/internal/models"
	"test-task/internal/queue"
	"test-task/internal/store"
)

func New(st store.Store, dl downloader.Downloader, cfg Config) *Service {
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 100
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 2
	}

	return &Service{
		st:    st,
		dl:    dl,
		tasks: map[string]*models.Task{},
		q:     queue.New(queue.Config{Size: cfg.QueueSize}),
	}
}

func (s *Service) Start(ctx context.Context, workers int) {
	if workers <= 0 {
		workers = 2
	}
	queue.Start(ctx, workers, s.q, s.runTask)
}

func (s *Service) Load() {
	if task, err := s.st.LoadTasks(); err != nil {
		s.tasks = task
	}
}
