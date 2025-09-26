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
	tasks, err := s.st.LoadTask()
	if err != nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tasks = tasks

	for _, t := range s.tasks {
		for i := range t.Files {
			if t.Files[i].Status == models.Running {
				t.Files[i].Status = models.Pending
			}
		}

		for _, file := range t.Files {
			if file.Status == models.Pending {
				t.Status = models.TaskPending
				s.q.Push(t.ID)
				break
			}
		}
	}
}
