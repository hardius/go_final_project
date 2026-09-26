package server

import (
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

func Run() error {
	api.Init()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := ":" + os.Getenv("TODO_PORT")
	if port == ":" {
		port = ":7540"
	}

	return http.ListenAndServe(port, nil)
}
