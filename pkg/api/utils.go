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
	_ = json.NewEncoder(w).Encode(data)
}

// checkDate валидирует и нормализует дату задачи согласно ожиданиям тестов addtask_4_test.go.
func checkDate(task *db.Task) error {
	now := time.Now()
	todayStr := now.Format(dateFormat)

	// Если дата пустая, берём сегодняшнюю.
	if task.Date == "" {
		task.Date = todayStr
	} else {
		// Проверяем корректность формата даты.
		t, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			return errors.New("дата представлена в формате, отличном от 20060102")
		}

		// Если дата в прошлом, поднимаем её до сегодняшней.
		if !afterNow(t, now) {
			task.Date = todayStr
		}
	}

	// Если есть правило повторения, просто валидируем его через NextDate,
	// но НЕ меняем исходную дату задачи.
	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			return err
		}
	}

	return nil
}
