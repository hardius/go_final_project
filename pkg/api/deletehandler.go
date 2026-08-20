package api

import (
	"database/sql"
	"errors"
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
		status := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			status = http.StatusBadRequest
		}

		writeJson(w, errWrap(err), status)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
