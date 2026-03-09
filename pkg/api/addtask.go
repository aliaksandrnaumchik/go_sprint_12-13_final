package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"todo-scheduler/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения тела запроса"})
		return
	}

	// Десериализуем JSON в структуру Task
	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Проверяем и корректируем дату
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем идентификатор созданной записи
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
