package api

import (
	"net/http"
	"todo-scheduler/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем задачу в формате JSON
	writeJSON(w, task)
}
