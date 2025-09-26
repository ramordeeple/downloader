package service

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"test-task/internal/models"
	"test-task/internal/util"
	"time"
)

type taskIDView struct{ id string }

func (v taskIDView) GetID() string { return v.id }

func (s *Service) NewTask(urls []string) string {
	seen := make(map[string]struct{}, len(urls))
	files := make([]models.File, 0, len(urls))

	for _, raw := range urls {
		u := strings.TrimSpace(raw)
		if u == "" {
			continue
		}

		pu, err := url.ParseRequestURI(u)
		if err != nil || pu.Scheme == "" || pu.Host == "" {
			continue
		}
		if _, dup := seen[u]; dup {
			continue
		}
		seen[u] = struct{}{}

		base := util.SanitizeFileName(path.Base(pu.Path))
		if base == "" || base == "/" || base == "." {
			base = time.Now().Format("20060102_150405")
		}

		files = append(files, models.File{
			URL:    u,
			Name:   base,
			Status: models.Pending,
		})
	}

	if len(files) == 0 {
		return ""
	}

	t := &models.Task{
		ID:        util.RandID(6),
		CreatedAt: time.Now(),
		Status:    models.TaskPending,
		Files:     files,
	}

	_ = os.MkdirAll(filepath.Join("data", t.ID), 0o755)

	s.mutex.Lock()
	if s.tasks == nil {
		s.tasks = make(map[string]*models.Task)
	}
	s.tasks[t.ID] = t
	_ = s.st.SaveTask(t)
	s.mutex.Unlock()

	s.q.Push(t.ID)
	return t.ID
}
