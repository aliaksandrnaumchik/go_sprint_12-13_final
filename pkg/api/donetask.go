package api

import (
	"net/http"
	"time"

	"todo-scheduler/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()

	// Если правило повторения отсутствует — удаляем задачу
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		// Если задача периодическая — вычисляем следующую дату
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		// Обновляем дату задачи в БД
		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	// Успешное выполнение — возвращаем пустой JSON-объект
	writeJSON(w, map[string]interface{}{})
}
