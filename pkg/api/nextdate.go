package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go1f/pkg/db"
)

const (
	DaysInWeek = 7 // Количество дней в неделе
)

var (
	MaxDaysInMonths = [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31} // Максимальное количество дней в месяце
	DayRule         = "d"                                                     // Индикатор дня в правиле задачи
	WeekRule        = "w"                                                     // Индикатор недели в правиле задачи
	MonthRule       = "m"                                                     // Индикатор месяца в правиле задачи
	YearRule        = "y"                                                     // Индикатор года в правиле задачи
)

// nextDayHandler обрабатывает GET-запрос по переданным в URL "date", "now" и "repeat" на возврат
// обновлённой соответствующей даты. Возвращает строку с датой в формате "20060102". В случае
// неудачи возвращает ошибку.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	nowString := r.URL.Query().Get("now")
	repeat := r.URL.Query().Get("repeat")
	var now time.Time
	if nowString == "" {
		now = time.Now().UTC()
	} else {
		var err error
		now, err = time.Parse(db.DateString, nowString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	nextDate, err := nextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprint(w, nextDate)
}

// nextDate возвращает строку с датой в формате "20060102" в соответствии с текущим временем (now),
// правилом (repeat) и датой старта (dstart) задачи. В случае неудачи возвращает пустую троку и ошибку.
func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", fmt.Errorf("правило повторения отсутствует")
	}
	date, err := time.Parse(db.DateString, dstart)
	if err != nil {
		return "", err
	}
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	rule := strings.Split(repeat, " ")
	if err := checkRule(rule[0], len(rule)); err != nil {
		return "", err
	}
	var nextDate time.Time
	switch rule[0] {

	// обработка правила дня
	case DayRule:
		dayList, err := parseStirngToArrInt(rule[1])
		if err != nil {
			return "", err
		}
		if err := checkDayRule(dayList); err != nil {
			return "", err
		}
		nextDate = nextDateDayRule(now, date, dayList[0])
	// обработка правила недели

	case WeekRule:
		weekList, err := parseStirngToArrInt(rule[1])
		if err != nil {
			return "", err
		}
		if err := checkWeekRule(weekList); err != nil {
			return "", err
		}
		nextDate = nextDateWeekRule(now, date, weekList)

	// обработка правила месяца
	case MonthRule:
		daysList, err := parseStirngToArrInt(rule[1])
		if err != nil {
			return "", err
		}
		var monthsList []int
		if len(rule) == 3 {
			monthsList, err = parseStirngToArrInt(rule[2])
			if err != nil {
				return "", err
			}
		} else {
			for i := 1; i < 13; i++ {
				monthsList = append(monthsList, i)
			}
		}
		if err := checkMonthRule(daysList, monthsList); err != nil {
			return "", err
		}
		nextDate = nextDateMonthRule(now, date, daysList, monthsList)

	// обработка правила года
	case YearRule:
		nextDate = nextDateYearRule(now, date)
	}
	return nextDate.Format(db.DateString), nil
}

// nextDateDayRule возвращает актуальную дату задачи для правила дня. Получает на вход
// текущее время (now), дату старта задачи (date) и шаг повторений (dayNumber).
func nextDateDayRule(now, date time.Time, dayNumber int) time.Time {
	nextDate := date
	if now.After(date) {
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, dayNumber)
		}
	} else {
		nextDate = date.AddDate(0, 0, dayNumber)
	}
	return nextDate
}

// nextDateWeekRule возвращает актуальную дату задачи для правила недели. Получает на вход
// текущее время (now), дату старта задачи (date) и слайс номеров дней недели (пн-1...вс-7)
// (dayNumber).
func nextDateWeekRule(now, date time.Time, dayList []int) time.Time {
	delta := DaysInWeek
	var startWeekDay int
	if now.After(date) {
		startWeekDay = int(now.Weekday())
	} else {
		startWeekDay = int(date.Weekday())
	}
	for _, v := range dayList {
		deltaV := v - startWeekDay
		if deltaV > 0 && deltaV < delta {
			delta = deltaV
		} else if (DaysInWeek + deltaV) < delta {
			delta = DaysInWeek + deltaV
		}
	}
	return now.AddDate(0, 0, delta)
}

// nextDateMonthRule возвращает актуальную дату задачи для правила месяца. Получает на вход
// текущее время (now), дату старта задачи (date) и слайсы номеров дней месяца (dayList) и
// месяцев (monthsList).
func nextDateMonthRule(now, date time.Time, dayList, monthsList []int) time.Time {
	var startDate time.Time
	var nextDate time.Time
	if now.After(date) {
		startDate = now
	} else {
		startDate = date
	}
	nextLeapYear := startDate.Year() + 1
	for !isLeapYear(nextLeapYear) {
		nextLeapYear++
	}
	var delta = 8 * 355 * 24 * time.Hour
	for _, m := range monthsList {
		for _, d := range dayList {
			if d > MaxDaysInMonths[m-1] {
				continue
			}
			switch {
			case m == 2 && d == 29:
				if isLeapYear(startDate.Year()) {
					nextDate = time.Date(startDate.Year(), time.Month(m), d, 0, 0, 0, 0, startDate.Location())
					if !nextDate.After(startDate) {
						nextDate = time.Date(nextLeapYear, time.Month(m), d, 0, 0, 0, 0, startDate.Location())
					}
				} else {
					nextDate = time.Date(nextLeapYear, time.Month(m), d, 0, 0, 0, 0, startDate.Location())
				}
			case m == 2 && (d == -1 || d == -2):
				if isLeapYear(startDate.Year()) {
					d = 29 + d + 1
					nextDate = time.Date(startDate.Year(), time.Month(m), d, 0, 0, 0, 0, startDate.Location())
					if !nextDate.After(startDate) {
						d = d - 1
						nextDate = time.Date(startDate.Year()+1, time.Month(m), d, 0, 0, 0, 0, startDate.Location())
					}
				} else {
					d = 28 + d + 1
					nextDate = time.Date(startDate.Year(), time.Month(m), d, 0, 0, 0, 0, startDate.Location())
					if !nextDate.After(startDate) {
						if isLeapYear(startDate.Year() + 1) {
							d = d + 1
						}
						nextDate = time.Date(startDate.Year()+1, time.Month(m), d, 0, 0, 0, 0, startDate.Location())
					}
				}
			default:
				if d == -1 || d == -2 {
					d = MaxDaysInMonths[m-1] + d + 1
				}
				nextDate = time.Date(startDate.Year(), time.Month(m), d, 0, 0, 0, 0, startDate.Location())
				if !nextDate.After(startDate) {
					nextDate = time.Date(startDate.Year()+1, time.Month(m), d, 0, 0, 0, 0, startDate.Location())
				}
			}
			deltaV := nextDate.Sub(startDate)
			if delta > deltaV {
				delta = deltaV
			}
		}
	}
	return startDate.Add(delta)
}

// nextDateYearRule возвращает актуальную дату задачи для правила года. Получает на вход
// текущее время (now) и дату старта задачи (date).
func nextDateYearRule(now, date time.Time) time.Time {
	var nextDate = date
	if now.After(date) {
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	} else {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate
}

// checkRule проверяет общие требования к правилу задачи. Принимает на вход индикатор правила (rule)
// и количество разделенных пробелами строк в правиле (lenRule). В случае несоответствия правила основным
// требованиям, возвращает ошибку.
func checkRule(rule string, lenRule int) error {
	if rule != DayRule && rule != WeekRule && rule != MonthRule && rule != YearRule {
		return fmt.Errorf("недопустимый символ, указывающий на тип правила: '%s' ('d', 'w', 'm' или 'y')", rule)
	}
	if (rule == DayRule || rule == WeekRule) && lenRule != 2 {
		return fmt.Errorf("после правила '%s' может быть описана только 1 группа", rule)
	}
	if rule == MonthRule && lenRule != 2 && lenRule != 3 {
		return fmt.Errorf("после правила '%s' могут быть описаны только 1 или 2 группы", rule)
	}
	if rule == YearRule && lenRule != 1 {
		return fmt.Errorf("после правила '%s' ничего больше быть не должно", rule)
	}
	return nil
}

// checkDayRule проверяет специфические требования к правилу дня. Принимает на вход слайс чисел,
// получившийся в результате конвертации второй из строк правила, разделенных пробелом (dayList).
// В случае несоответствия правила, возвращает ошибку.
func checkDayRule(dayList []int) error {
	if len(dayList) != 1 {
		return fmt.Errorf("избыточное количество параметров в правиле 'd'")
	}
	if dayList[0] < 1 || dayList[0] > 400 {
		return fmt.Errorf("недопустимое значение параметра в правиле 'd' (допускается от 1 до 400)")
	}
	return nil
}

// checkWeekRule проверяет специфические требования к правилу недели. Принимает на вход слайс чисел,
// получившийся в результате конвертации второй из строк правила, разделенных пробелом (weekList).
// В случае несоответствия правила, возвращает ошибку.
func checkWeekRule(weekList []int) error {
	for _, d := range weekList {
		if d < 1 || d > 7 {
			return fmt.Errorf("недопустимое число дня недели: %d!(допускается от 1 до 7)", d)
		}
	}
	return nil
}

// checkMonthRule проверяет специфические требования к правилу месяца. Принимает на вход два слайса чисел,
// получившихся в результате конвертации второй и третьей из строк правила, разделенных пробелом (daysList
// и monthsList соответственно). В случае несоответствия правила, возвращает ошибку.
func checkMonthRule(daysList, monthsList []int) error {
	for _, d := range daysList {
		if !((d >= 1 && d <= 31) || d == -1 || d == -2) {
			return fmt.Errorf("недопустимое число дня месяца: %d!(допускается от 1 до 31, -1, -2)", d)
		}
	}
	for _, m := range monthsList {
		if m < 1 || m > 12 {
			return fmt.Errorf("недопустимое число месяца: %d!(допускается от 1 до 12)", m)
		}
	}
	return nil
}

// parseStirngToArrInt конвертирует строку s в слайс чисел. В случае неудачи, возвращает
// пустой слайс и ошибку конвертации.
func parseStirngToArrInt(s string) ([]int, error) {
	sArr := strings.Split(s, ",")
	var sArrInt []int
	for _, v := range sArr {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return sArrInt, err
		}
		sArrInt = append(sArrInt, vInt)
	}
	return sArrInt, nil
}

// isLeapYear проверяте, является ли год year высокосным. Если да, возвращает true,
// иначе false.
func isLeapYear(year int) bool {
	if year%4 == 0 {
		if year%100 == 0 {
			return year%400 == 0
		}
		return true
	}
	return false
}
