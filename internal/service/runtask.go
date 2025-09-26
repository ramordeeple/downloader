package service

import (
	"context"
	"os"
	"path/filepath"
	"test-task/internal/downloader"
	"test-task/internal/models"
	"time"
)

func (s *Service) get(id string) *models.Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.tasks[id]
}

func (s *Service) runTask(ctx context.Context, id string) {
	s.mutex.Lock()
	t := s.tasks[id]
	s.mutex.Unlock()

	outDir := filepath.Join("data", t.ID)
	_ = os.MkdirAll(outDir, 0o755)

	dl := downloader.NewHTTPDownloader(downloader.Config{
		DownloadDir:   outDir,
		ClientTimeout: 60 * time.Second,
		MaxRetries:    2,
	})

	t.Status = models.TaskRunning
	_ = s.st.SaveTask(t)

	for i := range t.Files {
		f := &t.Files[i]
		f.Status = models.Running
		_ = s.st.SaveTask(t)

		name, size, err := dl.Fetch(context.Background(), f.URL, f.Name)
		if err != nil {
			f.Status, f.Error = models.Failed, err.Error()
		} else {
			f.Status, f.Name, f.SizeBytes = models.Completed, name, size
		}
		s.finishStatus(t)
	}

	s.finishStatus(t)
}
