package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"todo-scheduler/pkg/db"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	todayStr := now.Format(dateFormat)

	if task.Date == "" {
		task.Date = todayStr
	} else {
		t, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			return errors.New("дата представлена в формате, отличном от 20060102")
		}

		if !afterNow(t, now) {
			task.Date = todayStr
		}
	}

	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			return err
		}
	}

	return nil
}
