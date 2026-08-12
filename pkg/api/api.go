package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

type ErrorJson struct {
	Message string `json:"error"`
}

func writeJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func errWrap(err error) ErrorJson {
	return ErrorJson{err.Error()}
}
