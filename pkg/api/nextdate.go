package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now) || date.Equal(now)
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
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(dateFormat), nil

	case "w", "m":
		// Пока возвращаем ошибку для неподдерживаемых форматов
		return "", errors.New("неподдерживаемый формат правила повторения")

	default:
		return "", errors.New("неизвестный тип правила повторения")
	}
}
