package usecase

import (
	"sync"
	"test-task/internal/config"
	"test-task/internal/domain"
)

type TaskService struct {
	cfg   config.Service
	repo  TaskRepo
	fetch FileFetcher
	q     Queue
	log   Logger
	clk   Clock

	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewTaskService(r TaskRepo, f FileFetcher, q Queue, l Logger, c Clock, cfg config.Service) *TaskService {
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 100
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 2
	}

	return &TaskService{
		cfg:   cfg,
		repo:  r,
		fetch: f,
		q:     q,
		log:   l,
		clk:   c,
		tasks: make(map[string]*domain.Task),
	}
}
