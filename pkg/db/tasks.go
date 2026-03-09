package db

import (
	"database/sql"
	"fmt"
)

// Tasks возвращает список ближайших задач (не более limit записей), отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к БД: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки БД: %w", err)
		}
		tasks = append(tasks, &task)
	}

	// Если задач нет, возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask возвращает задачу по указанному ID или ошибку, если задача не найдена
func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	var task Task
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка выполнения запроса к БД: %w", err)
	}

	return &task, nil
}

// UpdateTask обновляет задачу в базе данных
func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи в БД: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества изменённых строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := DB.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи в БД: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества изменённых строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
