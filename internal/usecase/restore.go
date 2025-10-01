package usecase

import "test-task/internal/domain"

func (s *TaskService) Restore() error {
	all, err := s.repo.LoadAll()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for id, t := range all {
		normalizeTask(t)
		s.tasks[id] = t
		_ = s.repo.SaveTask(t)

		if t.Status != domain.TaskCompleted {
			s.q.Push(id)
		}
	}
	return nil
}

func normalizeTask(t *domain.Task) {
	completed, failed := 0, 0

	for i := range t.Files {
		switch t.Files[i].Status {
		case domain.Completed:
			completed++
		case domain.Failed:
			failed++
		default:
			t.Files[i].Status = domain.Pending
			t.Files[i].Error = ""
		}
	}

	switch {
	case completed == len(t.Files):
		t.Status = domain.TaskCompleted
	case failed > 0:
		t.Status = domain.TaskPending
	default:
		t.Status = domain.TaskPending
	}
}
