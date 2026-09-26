package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/signin", signInHandler)
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

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
				fmt.Println(jwt)
			}
			var valid bool
			valid = ValidateToken(jwt)
			fmt.Println(valid)
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
