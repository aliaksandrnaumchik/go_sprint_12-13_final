package api

import (
	"net/http"
	"todo-scheduler/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Удаляем задачу из БД
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешное удаление — возвращаем пустой JSON-объект
	writeJSON(w, map[string]interface{}{})
}
