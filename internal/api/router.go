package api

import "net/http"

type API struct {
	Tasks Tasks
}

func New(s TaskService) *API {
	return &API{
		Tasks: Tasks{S: s}}
}

func NewMux(a *API) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/tasks", a.Tasks.HandleCreate) // POST
	mux.HandleFunc("/tasks", a.Tasks.HandleGet)    // GET by id

	return mux
}
