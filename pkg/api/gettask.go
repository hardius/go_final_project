package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	writeJson(w, task, http.StatusOK)
}
