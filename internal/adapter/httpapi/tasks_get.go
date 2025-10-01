package httpapi

import (
	"net/http"
	"strings"
	"test-task/internal/util"
)

func (h Tasks) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		util.ErrorJSON(w, http.StatusBadRequest, "missing id")
		return
	}

	task := h.S.GetTask(id)
	if task == nil {
		util.ErrorJSON(w, http.StatusNotFound, "task not found")
		return
	}

	util.JSON(w, http.StatusOK, task)
}
