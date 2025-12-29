package httpapi

import (
	"downloader/internal/domain"
	"net/http"
)

type TaskUsecase interface {
	NewTask(urls []string) (string, error)
	GetTask(id string) (*domain.Task, error)
}

func New(uc TaskUsecase) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/tasks", &tasksCreate{uc: uc})
	mux.Handle("/tasks/", &tasksGet{uc: uc})
	return mux
}
