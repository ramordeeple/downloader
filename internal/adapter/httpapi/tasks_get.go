package httpapi

import (
	"net/http"
	"strings"
)

// TaskView — DTO для выдачи наружу (без лишних полей domain.Task).
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

type tasksGet struct {
	uc TaskUsecase
}

func (h *tasksGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, "missing task id", http.StatusBadRequest)
		return
	}

	t, err := h.uc.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, t)
}
