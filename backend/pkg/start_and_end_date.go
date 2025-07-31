package utils

import "time"

func GetStartAndEndDateVNTime(month, year int) (time.Time, time.Time) {
	locVN, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	startVN := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, locVN)
	firstOfNextMonthVN := startVN.AddDate(0, 1, 0)
	endVN := firstOfNextMonthVN.AddDate(0, 0, -1)

	return startVN, endVN
}
