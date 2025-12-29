package httpapi

import (
	"downloader/internal/domain"
	"time"
)

type TaskView struct {
	ID        string     `json:"id"`
	CreatedAt string     `json:"created_at"`
	Status    string     `json:"status"`
	Files     []FileView `json:"files"`
}

type FileView struct {
	URL       string `json:"url"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
}

func toView(t *domain.Task) *TaskView {
	v := &TaskView{
		ID:        t.ID,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
		Status:    string(t.Status),
		Files:     make([]FileView, 0, len(t.Files)),
	}
	for _, f := range t.Files {
		v.Files = append(v.Files, FileView{
			URL:       f.URL,
			Name:      f.Name,
			Status:    string(f.Status),
			Error:     f.Error,
			SizeBytes: f.SizeBytes,
		})
	}
	return v
}
