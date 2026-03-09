package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	// Сравниваем только даты, игнорируя время
	dateOnly := date.Truncate(24 * time.Hour)
	nowOnly := now.Truncate(24 * time.Hour)
	return dateOnly.After(nowOnly) || dateOnly.Equal(nowOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка пустой строки repeat
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Парсинг исходной даты
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("некорректный формат даты dstart")
	}

	parts := strings.Split(repeat, " ")
	ruleType := parts[0]

	switch ruleType {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d: требуется указать интервал в днях")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("некорректное число дней в правиле d")
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("интервал в днях должен быть от 1 до 400")
		}

		date := start
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(dateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y: не должно быть дополнительных параметров")
		}

		date := start
		// Проверяем, что дата валидна (например, 29 февраля в невисокосном году)
		if !isValidDate(date) {
			// Если дата невалидна, сдвигаем на следующий день
			date = date.AddDate(0, 0, 1)
		}

		for !afterNow(date, now) {
			nextYear := date.AddDate(1, 0, 0)
			if isValidDate(nextYear) {
				date = nextYear
			} else {
				// Если в следующем году дата невалидна (29.02), сдвигаем на 01.03
				date = nextYear.AddDate(0, 0, 1)
			}
		}
		return date.Format(dateFormat), nil

	default:
		return "", errors.New("неизвестный тип правила повторения")
	}
}

// isValidDate проверяет, является ли дата валидной (например, 29.02 в високосном году)
func isValidDate(t time.Time) bool {
	year, month, day := t.Date()
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location()).Day()
	return day <= lastDay
}
