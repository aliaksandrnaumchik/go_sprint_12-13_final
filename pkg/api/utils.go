package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"todo-scheduler/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата пустая, берём сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
		return nil
	}

	// Проверяем корректность формата даты
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("дата представлена в формате, отличном от 20060102")
	}

	// Если есть правило повторения, проверяем его и получаем следующую дату
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
		return nil
	}

	// Если дата в прошлом и нет правила повторения, берём сегодняшнюю дату
	if !afterNow(t, now) {
		task.Date = now.Format(dateFormat)
	}

	return nil
}
