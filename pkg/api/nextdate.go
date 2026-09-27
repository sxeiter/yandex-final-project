package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата начала (dstart): %w", err)
	}

	parts := strings.Split(strings.TrimSpace(repeat), " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("недопустимый интервал дней (должен быть от 1 до 400)")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила w")
		}
		daysStr := strings.Split(parts[1], ",")
		validDays := [8]bool{}
		for _, dStr := range daysStr {
			d, err := strconv.Atoi(strings.TrimSpace(dStr))
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("недопустимый день недели: %s", dStr)
			}
			validDays[d] = true
		}

		for {
			date = date.AddDate(0, 0, 1)
			goDay := int(date.Weekday())
			ruleDay := goDay
			if goDay == 0 {
				ruleDay = 7
			}
			if validDays[ruleDay] && date.After(now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("неверный формат правила m")
		}

		daysStr := strings.Split(parts[1], ",")
		validDays := [32]bool{}
		hasMinus1, hasMinus2 := false, false

		for _, dStr := range daysStr {
			dStr = strings.TrimSpace(dStr)
			switch dStr {
			case "-1":
				hasMinus1 = true
			case "-2":
				hasMinus2 = true
			default:
				d, err := strconv.Atoi(dStr)
				if err != nil || d < 1 || d > 31 {
					return "", fmt.Errorf("недопустимый день месяца: %s", dStr)
				}
				validDays[d] = true
			}
		}

		validMonths := [13]bool{}
		hasMonths := false
		if len(parts) == 3 {
			hasMonths = true
			monthsStr := strings.Split(parts[2], ",")
			for _, mStr := range monthsStr {
				m, err := strconv.Atoi(strings.TrimSpace(mStr))
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("недопустимый месяц: %s", mStr)
				}
				validMonths[m] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			monthOk := !hasMonths || validMonths[int(date.Month())]
			if !monthOk {
				continue
			}

			dayOk := validDays[date.Day()]
			if !dayOk && (hasMinus1 || hasMinus2) {
				lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
				if hasMinus1 && date.Day() == lastDay {
					dayOk = true
				}
				if hasMinus2 && date.Day() == lastDay-1 {
					dayOk = true
				}
			}

			if dayOk && date.After(now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", rule)
	}

	return date.Format(DateFormat), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Некорректный параметр now", http.StatusBadRequest)
			return
		}
	}

	nextDateStr, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDateStr))
}
