package service

import (
	"test-task/internal/models"
	"test-task/internal/util"
	"time"
)

type taskIDView struct{ id string }

func (v taskIDView) GetID() string { return v.id }

func (s *Service) NewTask(urls []string) string {
	files := make([]models.File, 0, len(urls))
	for _, u := range urls {
		files = append(files, models.File{URL: u, Status: models.Pending})
	}

	t := &models.Task{
		ID:        util.RandID(6),
		CreatedAt: time.Now(),
		Status:    models.TaskPending,
		Files:     files,
	}

	s.mutex.Lock()
	s.tasks[t.ID] = t
	_ = s.st.SaveTask(t)
	s.mutex.Unlock()

	s.q.Push(t.ID)
	return t.ID
}
