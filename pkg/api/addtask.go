package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.ccom/hardius/go_final_project/pkg/db"
)

var counter int = 1

type ErrorJson struct {
	Message string `json:"error"`
}

type IdJson struct {
	ID string `json:"id"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	fmt.Println(counter) // для отслеживания
	counter++

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		err = errors.New("Task title is empty")

		writeJson(w, errWrap(err), http.StatusBadRequest)

		fmt.Println("empty title") // для отслеживания
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		fmt.Println("чек дейт?", err)
		return
	}

	fmt.Println(task)
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	idJson := IdJson{ID: strconv.Itoa(int(id))}
	writeJson(w, idJson, http.StatusOK)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(layout)
	}

	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(layout)
		} else {
			task.Date = next
		}
	}

	return nil
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
