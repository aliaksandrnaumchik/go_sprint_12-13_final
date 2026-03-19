package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateOnly := date.Truncate(24 * time.Hour)
	nowOnly := now.Truncate(24 * time.Hour)
	return dateOnly.After(nowOnly) || dateOnly.Equal(nowOnly)
}

func lastDayOfMonth(year int, month time.Month, loc *time.Location) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, loc).Day()
}

func addYearRule(t time.Time) time.Time {
	year, month, day := t.Date()
	loc := t.Location()

	if month == time.February && day == 29 {
		return time.Date(year+1, time.March, 1, 0, 0, 0, 0, loc)
	}

	return t.AddDate(1, 0, 0)
}

func parseMonthDays(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("некорректный день месяца")
		}
		if v == 0 || v < -2 || v > 31 {
			return nil, errors.New("некорректное значение дня месяца")
		}
		result = append(result, v)
	}
	if len(result) == 0 {
		return nil, errors.New("список дней месяца пуст")
	}
	return result, nil
}

func parseMonths(s string) (map[int]bool, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	result := make(map[int]bool, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("некорректный месяц")
		}
		if v < 1 || v > 12 {
			return nil, errors.New("месяц должен быть в диапазоне 1-12")
		}
		result[v] = true
	}
	if len(result) == 0 {
		return nil, errors.New("список месяцев пуст")
	}
	return result, nil
}

func nextMonthlyOccurrence(prev time.Time, days []int, months map[int]bool) time.Time {
	loc := prev.Location()
	year, month, day := prev.Date()

	for {
		if months != nil && !months[int(month)] {
			month++
			if month > 12 {
				month = 1
				year++
			}
			day = 0
			continue
		}

		lastDay := lastDayOfMonth(year, month, loc)
		var candidate *time.Time

		for _, d := range days {
			actualDay := d
			if d < 0 {
				actualDay = lastDay + d + 1
			}
			if actualDay < 1 || actualDay > lastDay {
				continue
			}

			if day > 0 && actualDay <= day {
				continue
			}

			t := time.Date(year, month, actualDay, 0, 0, 0, 0, loc)
			if candidate == nil || t.Before(*candidate) {
				tmp := t
				candidate = &tmp
			}
		}

		if candidate != nil {
			return *candidate
		}

		month++
		if month > 12 {
			month = 1
			year++
		}
		day = 0
	}
}

func parseWeekdays(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("некорректный день недели")
		}
		if v < 1 || v > 7 {
			return nil, errors.New("день недели должен быть в диапазоне 1-7")
		}
		result = append(result, v)
	}
	if len(result) == 0 {
		return nil, errors.New("список дней недели пуст")
	}
	return result, nil
}

func weekdayNumber(t time.Time) int {
	return ((int(t.Weekday()) + 6) % 7) + 1
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("некорректный формат даты dstart")
	}

	fields := strings.Fields(repeat)
	if len(fields) == 0 {
		return "", errors.New("пустое правило повторения")
	}

	ruleType := fields[0]

	switch ruleType {
	case "d":
		if len(fields) != 2 {
			return "", errors.New("неверный формат правила d")
		}
		interval, err := strconv.Atoi(fields[1])
		if err != nil {
			return "", errors.New("некорректное число дней в правиле d")
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("интервал в днях должен быть от 1 до 400")
		}

		date := start.AddDate(0, 0, interval)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(dateFormat), nil

	case "y":
		if len(fields) != 1 {
			return "", errors.New("неверный формат правила y")
		}

		date := addYearRule(start)
		for !afterNow(date, now) {
			date = addYearRule(date)
		}
		return date.Format(dateFormat), nil

	case "m":
		if len(fields) < 2 || len(fields) > 3 {
			return "", errors.New("неверный формат правила m")
		}

		days, err := parseMonthDays(fields[1])
		if err != nil {
			return "", err
		}

		var months map[int]bool
		if len(fields) == 3 {
			months, err = parseMonths(fields[2])
			if err != nil {
				return "", err
			}
		}

		date := nextMonthlyOccurrence(start, days, months)
		for !afterNow(date, now) {
			date = nextMonthlyOccurrence(date, days, months)
		}
		return date.Format(dateFormat), nil

	case "w":
		if len(fields) != 2 {
			return "", errors.New("неверный формат правила w")
		}

		weekdays, err := parseWeekdays(fields[1])
		if err != nil {
			return "", err
		}
		weekSet := make(map[int]bool, len(weekdays))
		for _, w := range weekdays {
			weekSet[w] = true
		}

		date := start.AddDate(0, 0, 1)
		for {
			if date.After(now) && weekSet[weekdayNumber(date)] {
				return date.Format(dateFormat), nil
			}
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", errors.New("неизвестный тип правила повторения")
	}
}
