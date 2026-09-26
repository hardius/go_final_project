package api

import (
	"errors"
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		err := errors.New("This method isn't handled.")
		writeJson(w, errWrap(err), http.StatusBadRequest)
	}
}
