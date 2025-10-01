package httpapi

import (
	"net/http"
)

type TaskUsecase interface {
	NewTask(urls []string) (string, error)
	GetTask(id string) (*TaskView, error)
}

func New(uc TaskUsecase) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/tasks", &tasksCreate{uc: uc})

	mux.Handle("/tasks/", &tasksGet{uc: uc})

	return mux
}
