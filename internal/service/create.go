package service

import (
	"test-task/internal/models"
	"test-task/internal/util"
	"time"
)

type taskIDView struct{ id string }

func (v taskIDView) GetID() string { return v.id }

func (s *Service) NewTask(urls []string) taskIDView {
	files := make([]models.File, 0, len(urls))
	for _, u := range urls {
		files = append(files, models.File{URL: u, Status: models.Pending})
	}

	task := &models.Task{
		ID:        util.RandID(6),
		CreatedAt: time.Now().String(),
		Status:    models.TaskPending,
		Files:     files,
	}

	s.mutex.Lock()
	s.tasks[task.ID] = task
	_ = s.st.SaveTask(task)
	s.mutex.Unlock()

	s.q.Push(task.ID)
	return taskIDView{id: task.ID}
}
