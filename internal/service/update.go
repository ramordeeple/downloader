package service

import "test-task/internal/models"

func (s *Service) updateTask(task *models.Task, status models.TaskStatus) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task.Status = status
}

func (s *Service) updateFile(task *models.Task, i int, status models.FileStatus, err string, size int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task.Files[i].Status = status
	task.Files[i].Error = err
	if size > 0 {
		task.Files[i].SizeBytes = size
	}
}

func (s *Service) setFileName(task *models.Task, i int, name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task.Files[i].Name = name
}

func (s *Service) finishStatus(task *models.Task) {
	ok, bad := 0, 0
	for _, file := range task.Files {
		switch file.Status {
		case models.Completed:
			ok++
		case models.Failed:
			bad++
		}
	}

	switch {
	case ok == len(task.Files):
		s.updateTask(task, models.TaskCompleted)
	case bad > 0:
		s.updateTask(task, models.TaskFailed)
	default:
		s.updateTask(task, models.TaskRunning)
	}
}
