package api

import (
	"io"
	"net/http"
	"time"

	"go_final_project/pkg/functions"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowString == "" {
		now = time.Now()
	} else {
		res, err := time.Parse(layout, nowString)
		if err != nil {
			http.Error(w, "Wrong format of date.", http.StatusBadRequest)
			return
		}

		now = res
	}

	result, err := functions.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	io.WriteString(w, result)
}
