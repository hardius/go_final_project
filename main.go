package main

import (
	"fmt"
	"os"

	"github.ccom/hardius/go_final_project/pkg/db"
	"github.ccom/hardius/go_final_project/pkg/server"
)

func main() {
	if errDB := db.Init(os.Getenv("TODO_DBFILE")); errDB != nil {
		fmt.Println(errDB)
	}

	if errServer := server.Run(); errServer != nil {
		fmt.Println(errServer)
	}
}
