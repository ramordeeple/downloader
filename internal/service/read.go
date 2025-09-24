package service

import "time"

type FileView struct {
	URL       string `json:"url"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Error     string `json:"error"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
}

type TaskView struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Status    string     `json:"status"`
	Files     []FileView `json:"files"`
}

func (s *Service) GetTask(id string) *TaskView {
	s.mutex.Lock()
	task := s.tasks[id]
	s.mutex.Unlock()
	if task == nil {
		return nil
	}
	out := &TaskView{
		ID:        task.ID,
		CreatedAt: task.CreatedAt,
		Status:    string(task.Status),
		Files:     make([]FileView, len(task.Files)),
	}

	for i, f := range task.Files {
		out.Files[i] = FileView{
			URL:       f.URL,
			Name:      f.Name,
			Status:    string(f.Status),
			Error:     f.Error,
			SizeBytes: f.SizeBytes,
		}
	}

	return out
}
