package api

import (
	"net/http"

	"github.com/hardius/go_final_project/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, errWrap(err), http.StatusInternalServerError)
			return
		}
	} else {
		err = db.UpdateDate(task)
		if err != nil {
			writeJson(w, errWrap(err), http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
