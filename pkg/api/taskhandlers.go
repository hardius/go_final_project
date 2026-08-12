package api

import "net/http"

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
