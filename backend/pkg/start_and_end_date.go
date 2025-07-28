package utils

import "time"

func GetStartAndEndDate(month, year int) (time.Time, time.Time) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	firstOfNextMonth := startDate.AddDate(0, 1, 0)

	endDate := firstOfNextMonth.AddDate(0, 0, -1)

	return startDate, endDate
}
