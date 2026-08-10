package server

import (
	"net/http"
	"os"

	"github.ccom/hardius/go_final_project/pkg/api"
)

func Run() error {
	api.Init()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := ":" + os.Getenv("TODO_PORT")

	return http.ListenAndServe(port, nil)
}
