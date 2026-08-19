package api

import (
	"net/http"
	"strconv"

	"go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
