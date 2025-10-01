package httpapi

import (
	"net/http"
	"test-task/internal/api"
	"test-task/internal/usecase"
)

type API struct {
	Tasks Tasks
}

func NewAPI(s *usecase.Service) *API {
	return &API{
		Tasks: Tasks{S: &api.serviceAdapter{s}}}
}

func NewMux(a *API) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/tasks", a.Tasks.HandleCreate) // POST
	mux.HandleFunc("/tasks/", a.Tasks.HandleGet)   // GET by id

	return mux
}
