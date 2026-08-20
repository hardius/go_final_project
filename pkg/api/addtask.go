package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/functions"
)

type IdJson struct {
	ID string `json:"id"`
}

const layout = "20060102"

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		err = errors.New("Task title is empty")

		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, errWrap(err), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, errWrap(err), http.StatusInternalServerError)
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
		next, err = functions.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

	}

	if functions.AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(layout)
		} else {
			task.Date = next
		}
	}

	return nil
}
