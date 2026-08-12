package main

import (
	"fmt"
	"os"

	"github.ccom/hardius/go_final_project/pkg/db"
	"github.ccom/hardius/go_final_project/pkg/server"
)

func main() {
	db, errDB := db.Init(os.Getenv("TODO_DBFILE"))
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	if errServer := server.Run(); errServer != nil {
		fmt.Println(errServer)
	}
}
