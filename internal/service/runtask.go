package service

import (
	"context"
	"test-task/internal/models"
)

func (s *Service) get(id string) *models.Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.tasks[id]
}

func (s *Service) runTask(ctx context.Context, id string) {
	task := s.get(id)
	if task == nil {
		return
	}

	s.updateTask(task, models.TaskRunning)
	for i := range task.Files {
		s.updateFile(task, i, models.Running, "", 0)

		saved, size, err := s.dl.Fetch(ctx, task.Files[i].URL, task.Files[i].Name)
		if err != nil {
			s.updateFile(task, i, models.Failed, err.Error(), size)
		} else {
			s.setFileName(task, i, saved)
			s.updateFile(task, i, models.Completed, "", size)
		}

		_ = s.st.SaveTask(task)
	}

	s.finishStatus(task)
	_ = s.st.SaveTask(task)
}
