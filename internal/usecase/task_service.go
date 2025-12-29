// internal/usecase/task_service.go
package usecase

import (
	"context"
	"downloader/internal/config"
	"path"
	"strings"
	"sync"
	"time"

	"downloader/internal/domain"
	"downloader/internal/platform/util"
)

type TaskRepo interface {
	SaveTask(*domain.Task) error
	LoadAll() (map[string]*domain.Task, error)
}

type FileFetcher interface {
	Fetch(ctx context.Context, url, suggestedName, outDir string) (string, int64, error)
}

type Queue interface {
	Push(id string)
	Pop() <-chan string
	Close()
}

type Extractor interface {
	Extract(ctx context.Context, url string) ([]Media, error)
}

type Clock interface{ Now() time.Time }

type TaskService struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task

	repo  TaskRepo
	fetch FileFetcher
	q     Queue
	clk   Clock

	cfg     config.Service
	dataDir string
}

func NewTaskService(repo TaskRepo, fetch FileFetcher, q Queue, clk Clock, dataDir string) *TaskService {
	return &TaskService{
		repo:    repo,
		fetch:   fetch,
		q:       q,
		clk:     clk,
		dataDir: dataDir,
		tasks:   make(map[string]*domain.Task),
	}
}

func (s *TaskService) NewTask(urls []string) (string, error) {
	files := make([]domain.File, 0, len(urls))
	for _, u := range urls {
		name := util.SanitizeFileName(path.Base(strings.TrimSpace(u)))
		if name == "" || name == "/" || name == "." {
			name = util.RandID(4)
		}
		files = append(files, domain.File{
			URL:    u,
			Name:   name,
			Status: domain.Pending,
		})
	}
	id := util.RandID(6)
	t := &domain.Task{
		ID:        id,
		CreatedAt: s.clk.Now(),
		Status:    domain.TaskPending,
		Files:     files,
	}

	s.mu.Lock()
	s.tasks[id] = t
	_ = s.repo.SaveTask(t)
	s.mu.Unlock()

	s.q.Push(id)
	return id, nil
}

func (s *TaskService) GetTask(id string) (*domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id], nil
}
