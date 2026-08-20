package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"go_final_project/pkg/db"
)

func putHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		err = errors.New("Task title is empty")

		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
