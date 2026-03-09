package api

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateMidnight := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateMidnight.After(nowMidnight) || dateMidnight.Equal(nowMidnight)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

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
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("интервал дней должен быть от 1 до 400")
		}

		date := start
		// Пока дата не будет после now, прибавляем интервал
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, days)
		}
		return date.Format(dateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y")
		}

		date := start
		// Если дата старта уже прошла относительно now, ищем следующий год
		if !afterNow(date, now) {
			// Начинаем с года после даты старта
			date = time.Date(start.Year()+1, start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		}

		// Ищем первый год, где дата будет после now
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(dateFormat), nil

	case "m":
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила m: требуется указать дни месяца")
		}

		dayStrs := strings.Split(parts[1], ",")
		var days []int
		for _, ds := range dayStrs {
			d, err := strconv.Atoi(ds)
			if err != nil {
				return "", errors.New("некорректные дни в правиле m")
			}
			days = append(days, d)
		}

		current := start.AddDate(0, 1, 0) // начинаем со следующего месяца
		maxIterations := 100              // защита от бесконечного цикла

		for i := 0; i < maxIterations; i++ {
			year, month, _ := current.Date()
			lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, current.Location()).Day()

			var candidates []time.Time
			for _, day := range days {
				targetDay := day
				if day < 0 {
					targetDay = lastDay + day + 1
				}
				if targetDay < 1 || targetDay > lastDay {
					continue // пропускаем некорректные дни
				}

				candidate := time.Date(year, month, targetDay, 0, 0, 0, 0, current.Location())
				candidates = append(candidates, candidate)
			}

			// Сортируем кандидаты по дате
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].Before(candidates[j])
			})

			// Проверяем кандидатов в порядке возрастания дат
			for _, candidate := range candidates {
				if afterNow(candidate, now) {
					return candidate.Format(dateFormat), nil
				}
			}

			// Переходим к следующему месяцу
			current = current.AddDate(0, 1, 0)
		}
		return "", errors.New("не удалось найти подходящую дату в течение 100 месяцев")

	case "w":
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила w: требуется указать дни недели")
		}

		dayStrs := strings.Split(parts[1], ",")
		var targetDays []int
		for _, ds := range dayStrs {
			d, err := strconv.Atoi(ds)
			if err != nil {
				return "", errors.New("некорректные дни недели в правиле w")
			}
			targetDays = append(targetDays, d)
		}

		current := start
		maxIterations := 365 // защита от бесконечного цикла (максимум год)

		for i := 0; i < maxIterations; i++ {
			current = current.AddDate(0, 0, 1)
			weekday := int(current.Weekday())

			for _, target := range targetDays {
				if weekday == target {
					if afterNow(current, now) {
						return current.Format(dateFormat), nil
					}
					break
				}
			}
		}
		return "", errors.New("не удалось найти подходящий день недели в течение года")

	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}
