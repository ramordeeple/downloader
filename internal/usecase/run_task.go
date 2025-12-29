package usecase

import (
	"context"
	"os"
	"path/filepath"

	"downloader/internal/domain"
)

func (s *TaskService) runTask(ctx context.Context, id string) {
	s.mu.RLock()
	t := s.tasks[id]
	s.mu.RUnlock()
	if t == nil {
		return
	}

	outDir := filepath.Join(s.dataDir, t.ID)
	_ = os.MkdirAll(outDir, 0o755)

	t.Status = domain.TaskRunning
	_ = s.repo.SaveTask(t)

	for i := range t.Files {
		f := &t.Files[i]
		if f.Status == domain.Completed {
			continue
		}

		f.Status = domain.Running
		_ = s.repo.SaveTask(t)

		name, size, err := s.fetch.Fetch(ctx, f.URL, f.Name, outDir)
		if err != nil {
			f.Status = domain.Failed
			f.Error = err.Error()
		} else {
			f.Status = domain.Completed
			f.Name = name
			f.SizeBytes = size
		}
		s.updateTaskStatus(t)
	}

	s.updateTaskStatus(t)
}

func (s *TaskService) updateTaskStatus(t *domain.Task) {
	completed, failed := 0, 0
	for _, f := range t.Files {
		switch f.Status {
		case domain.Completed:
			completed++
		case domain.Failed:
			failed++
		}
	}

	switch {
	case completed == len(t.Files):
		t.Status = domain.TaskCompleted
	case failed > 0:
		t.Status = domain.TaskFailed
	default:
		t.Status = domain.TaskRunning
	}
	_ = s.repo.SaveTask(t)
}
