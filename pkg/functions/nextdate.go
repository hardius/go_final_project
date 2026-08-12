package functions

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.ccom/hardius/go_final_project/pkg/db"
)

const initialMin = 100
const layout = "20060102"

func afterNow(date, now time.Time) bool {
	return date.Format(layout) > now.Format(layout)
}

func lastDay(month time.Month, year int) int {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Now().Location()).Day()

	return lastDay
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Repeat string is empty.")
	}

	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", errors.New("Wrong date.")
	}

	partRepeatStr := strings.Split(repeat, " ")

	var fresult func(time.Time, time.Time, []string) (string, error)

	switch partRepeatStr[0] {
	case "d":
		fresult = ruleD
	case "y":
		fresult = ruleY
	case "w":
		fresult = ruleW
	case "m":
		fresult = ruleM
	default:
		return "", errors.New("Unknown command.")
	}

	result, err := fresult(now, date, partRepeatStr)
	return result, err
}

func ruleD(now, dstart time.Time, partRepeatStr []string) (string, error) {
	if len(partRepeatStr) != 2 {
		return "", errors.New("Invalid number of arguments.")
	}

	val, err := strconv.Atoi(partRepeatStr[1])
	if err != nil {
		return "", err
	}

	if val > 400 {
		return "", errors.New("Argument for d-rule cannot be over 400.")
	}

	for {
		dstart = dstart.AddDate(0, 0, val)
		if afterNow(dstart, now) {
			result := dstart.Format(layout)
			return result, nil
		}
	}
}

func ruleY(now, dstart time.Time, partRepeatStr []string) (string, error) {
	if len(partRepeatStr) > 1 {
		return "", errors.New("y-rule doesn't need any argument")
	}

	var date time.Time
	if afterNow(dstart, now) {
		result := dstart.AddDate(1, 0, 0).Format(layout)
		return result, nil
	}

	date = dstart.AddDate(now.Year()-dstart.Year(), 0, 0)
	if afterNow(date, now) {
		result := date.Format(layout)
		return result, nil
	}

	result := date.AddDate(1, 0, 0).Format(layout)
	return result, nil
}

func ruleW(now, dstart time.Time, partRepeatStr []string) (string, error) {
	if len(partRepeatStr) != 2 {
		return "", errors.New("Wrong use of ruleW.")
	}

	var date time.Time
	if afterNow(dstart, now) {
		date = dstart
	} else {
		date = now
	}

	dateWD := int(date.Weekday())
	if dateWD == 0 {
		dateWD = 7
	}

	arguments := strings.Split(partRepeatStr[1], ",")

	var min int = initialMin
	var diff int
	for _, v := range arguments {
		currWD, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}

		if currWD > 7 {
			return "", errors.New("Argument for w-rule cannot be over 7")
		}

		if currWD <= dateWD {
			currWD = currWD + 7
		}

		diff = currWD - dateWD
		if diff < min {
			min = diff
		}
	}

	result := date.AddDate(0, 0, min).Format(layout)
	return result, nil
}

func ruleM(now, dstart time.Time, partRepeatStr []string) (string, error) {
	if len(partRepeatStr) > 3 && len(partRepeatStr) < 2 {
		return "", errors.New("Wrong use of ruleWM.")
	}

	var result string

	var date time.Time
	if afterNow(dstart, now) {
		date = dstart
	} else {
		date = now
	}

	arguments1 := strings.Split(partRepeatStr[1], ",")

	day := date.Day()
	month := int(date.Month())

	if len(partRepeatStr) > 2 {

		arguments2 := strings.Split(partRepeatStr[2], ",")

		var mindiffMonth int = initialMin
		var diff int

		var checkCurrMonth bool
		for _, v := range arguments2 {
			currMonth, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}

			if currMonth > 12 {
				return "", errors.New("Month arguments for m-rule cannot be over 12")
			}
			if currMonth > month {
				diff = currMonth - month
			} else if currMonth < month {
				diff = currMonth + 12 - month
			} else if currMonth == month {
				checkCurrMonth = true
				continue
			}

			if mindiffMonth > diff {
				mindiffMonth = diff
			}
		}

		var foundNext bool
		var mindiffDay int = initialMin

		if checkCurrMonth {
			lastDay := lastDay(time.Month(month), date.Year())
			for _, v := range arguments1 {
				currDay, err := strconv.Atoi(v)
				if err != nil {
					return "", err
				}

				if currDay > 31 {
					return "", errors.New("Day arguments for m-rule cannot be over 31")
				}

				if currDay < 0 {
					currDay = lastDay + 1 + currDay
				}

				if currDay > day {
					foundNext = true
					diff = currDay - day
					if diff < mindiffDay {
						mindiffDay = diff
					}
				}
			}
		}

		if foundNext {
			result = date.AddDate(0, 0, mindiffDay).Format(layout)
		} else {
			date = date.AddDate(0, mindiffMonth, 0)
			month = month + mindiffMonth

			lastDay := lastDay(time.Month(month), date.Year())
			var minDay int = initialMin
			for _, v := range arguments1 {
				currDay, err := strconv.Atoi(v)
				if err != nil {
					return "", err
				}

				if currDay > 31 {
					return "", errors.New("Number of day cannot be over 31.")
				}

				if currDay < 0 {
					currDay = lastDay + 1 + currDay
				}

				if currDay < minDay {
					minDay = currDay
				}
			}

			date = time.Date(date.Year(), date.Month(), minDay,
				date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
			result = date.Format(layout)
		}
	} else {
		var mindiffDay int = initialMin
		var diff int
		lastDay := lastDay(time.Month(month), date.Year())
		for _, v := range arguments1 {
			currDay, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}

			if currDay > 31 {
				return "", errors.New("Number of day cannot be over 31.")
			}

			if currDay < 0 {
				currDay = lastDay + 1 + currDay
			} else if currDay == 31 && lastDay != 31 {
				currDay = currDay + lastDay
			}

			if currDay <= day {
				currDay = currDay + lastDay
			}

			diff = currDay - day
			if diff < mindiffDay {
				mindiffDay = diff
			}
		}

		result = date.AddDate(0, 0, mindiffDay).Format(layout)
	}

	return result, nil
}

func CheckDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(layout)
	}

	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(layout)
		} else {
			task.Date = next
		}
	}

	return nil
}
