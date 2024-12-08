package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// возвращает следующее значение даты задачи в соответствии с правилом повторения
func NextDate(now time.Time, date string, repeat string) (NewDate string, Err error) {
	//разделяем правило по пробелам
	rulesList := strings.Fields(repeat)

	//проверяем парсится ли date
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return NewDate, fmt.Errorf("failed to convert date")
	}
	//проверяем длину rulesList
	if len(rulesList) < 1 || len(rulesList) > 3 {
		return NewDate, fmt.Errorf("invalid repeat rules format")
	}

	//проверяем первую букву
	switch rulesList[0] {
	case "y":
		if len(rulesList) == 1 {
			diffYears := 0
			parsedDate = parsedDate.AddDate(1, 0, 0)
			if parsedDate.Compare(now) == -1 {
				diffYears = now.Year() - parsedDate.Year()
			}
			parsedDate = parsedDate.AddDate(diffYears, 0, 0)
			NewDate = parsedDate.Format("20060102")
			return NewDate, nil
		}
	case "d":
		if len(rulesList) == 2 {
			nDays, err := strconv.Atoi(rulesList[1])
			if err != nil || nDays < 1 || nDays > 400 {
				return NewDate, fmt.Errorf("неверный формат количества дней в правиле повторения задачи")
			}
			changeFlag := false
			for parsedDate.Compare(now) == -1 {
				parsedDate = parsedDate.AddDate(0, 0, nDays)
				changeFlag = true
			}
			if !changeFlag {
				parsedDate = parsedDate.AddDate(0, 0, nDays)
			}
			NewDate = parsedDate.Format("20060102")
			return NewDate, nil
		}
	case "w":
		if len(rulesList) == 2 {
			var dayToCompare = int(parsedDate.Weekday())
			var rDate = parsedDate
			minDiff := 7

			if parsedDate.Compare(now) == -1 {
				dayToCompare = int(now.Weekday())
				rDate = now
			}

			if dayToCompare == 0 {
				dayToCompare += 7 //.Weekday() для воскресенья возвращает 0
			}

			weekdaysList := strings.Split(rulesList[1], ",")
			for _, w := range weekdaysList {
				numDay, err := strconv.Atoi(w)
				if err != nil || numDay < 1 || numDay > 7 {
					return NewDate, fmt.Errorf("неверный формат дней недели в правиле повторения задачи")
				}
				//ищем самый близкий день недели
				diff := numDay - dayToCompare
				if diff <= 0 {
					diff += 7 //добавляем неделю
				}
				if diff < minDiff {
					minDiff = diff
				}

			}
			NewDate = rDate.AddDate(0, 0, minDiff).Format("20060102")
			return NewDate, nil
		}
	case "m":
		if len(rulesList) > 1 {
			var rDate = parsedDate

			if parsedDate.Compare(now) == -1 {
				rDate = now
			}

			monthToCompare := int(rDate.Month())
			var monthNums, sortedByMonth, days []int
			var firstDates []time.Time
			var index int

			//смотрим какие месяца нам нужны
			if len(rulesList) == 3 {
				monthsList := strings.Split(rulesList[2], ",")
				for _, m := range monthsList {
					numMonth, err := strconv.Atoi(m)
					if err != nil || numMonth < 1 || numMonth > 12 {
						return NewDate, fmt.Errorf("неверный формат месяцев в правиле повторения задачи")
					}
					monthNums = append(monthNums, numMonth)
				}
			} else {
				//если месяца не указаны, будет работать с двумя ближайщими
				monthNums = append(monthNums, monthToCompare, int(rDate.AddDate(0, 1, 0).Month()))
			}
			//сортируем месяца в порядке возрастания
			monthNums = quickSort(monthNums)

			//запоминем индекс месяца, ближайщий к имеющейся дате
			for i, v := range monthNums {
				if v < monthToCompare {
					index = i + 1
				}
			}
			// теперь первым в списке будет текущий месяц (нужно для поиска ближайшей даты)
			sortedByMonth = append(sortedByMonth, monthNums[index:]...)
			sortedByMonth = append(sortedByMonth, monthNums[:index]...)

			for i := range sortedByMonth {
				//смотрим как далеко месяцы от monthToCompare
				sortedByMonth[i] = sortedByMonth[i] - monthToCompare
				if sortedByMonth[i] < 0 {
					sortedByMonth[i] += 12
				}
				dateToAdd := rDate.AddDate(0, sortedByMonth[i], -rDate.Day()+1)
				//массив с датами первых чисел месяцов
				firstDates = append(firstDates, dateToAdd)
			}

			//цикл по дням
			daysList := strings.Split(rulesList[1], ",")
			for _, d := range daysList {
				day, err := strconv.Atoi(d)
				if err != nil || day < -2 || day > 31 || day == 0 {
					return NewDate, fmt.Errorf("неверный формат дней в правиле повторения задачи")
				}
				days = append(days, day)
			}

			days = quickSort(days)
			for _, v := range firstDates {
				localDays := make([]int, len(days))
				copy(localDays, days)
				for i := 0; i < len(localDays); {
					if localDays[i] == -1 {
						localDays[i] = int(v.AddDate(0, 1, 0).Sub(v).Hours() / 24)
						localDays = quickSort(localDays)
						i = 0
						continue
					} else if localDays[i] == -2 {
						localDays[i] = int(v.AddDate(0, 1, 0).Sub(v).Hours()/24) - 1
						localDays = quickSort(localDays)
						i = 0
						continue
					}
					midDate := v.AddDate(0, 0, localDays[i]-1)
					if midDate.Compare(rDate) != 1 || v.Month() != midDate.Month() {
						i++
						continue
					}
					NewDate = midDate.Format("20060102")
					return NewDate, nil
				}
			}

		}
	}
	err = fmt.Errorf("неверный формат правила повторения задачи")
	return NewDate, err
}

// быстрая сортировка
func quickSort(nums []int) []int {
	if len(nums) < 2 {
		return nums
	}

	var less, greater []int
	mid := nums[len(nums)/2]

	for i := range nums {
		if i == len(nums)/2 {
			continue
		}
		if nums[i] <= mid {
			less = append(less, nums[i])
		} else {
			greater = append(greater, nums[i])
		}
	}
	qsLeft := quickSort(less)
	qsRight := quickSort(greater)
	return append(append(qsLeft, mid), qsRight...)
}
