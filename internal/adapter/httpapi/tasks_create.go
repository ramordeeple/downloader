package httpapi

import (
	"encoding/json"
	"net/http"
	"test-task/internal/usecase"
	"test-task/internal/util"
)

type TaskID interface {
	GetID() string
}

type serviceAdapter struct {
	s *usecase.Service
}

func (sa *serviceAdapter) NewTask(urls []string) string {
	return sa.s.NewTask(urls)
}

func (sa *serviceAdapter) GetTask(id string) any {
	return sa.s.GetTask(id)
}

type TaskService interface {
	NewTask(urls []string) string
	GetTask(id string) any
}

type Tasks struct {
	S TaskService
}

func (h Tasks) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		URLs []string `json:"urls"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.URLs) == 0 {
		util.ErrorJSON(w, http.StatusBadRequest, "body is empty")
		return
	}

	id := h.S.NewTask(body.URLs)
	util.JSON(w, http.StatusAccepted, map[string]string{"id": id})
}
