package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetCurrentDate() time.Time {
	return time.Now()
}

func GetCurrentTimeHCMCity() (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Errorf("Không thể load location: %w", err.Error())
	}
	timeNow := time.Now().In(loc)
	return timeNow, err
}

func GenerateCode(prefix string, digits int, getLastCodeFunc func() (string, error)) (string, error) {
	lastCode, err := getLastCodeFunc()
	if err != nil {
		if err.Error() == "record not found" {
			return fmt.Sprintf("%s%0*d", prefix, digits, 1), nil
		}
		return "", fmt.Errorf("failed to get last code: %w", err)
	}

	if lastCode == "" {
		return fmt.Sprintf("%s%0*d", prefix, digits, 1), nil
	}

	if !strings.HasPrefix(lastCode, prefix) {
		return "", fmt.Errorf("invalid code format: %s, expected prefix: %s", lastCode, prefix)
	}

	numberStr := strings.TrimPrefix(lastCode, prefix)
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse code number: %w", err)
	}

	nextNumber := number + 1

	maxNumber := int(math.Pow10(digits)) - 1
	if nextNumber > maxNumber {
		return "", fmt.Errorf("maximum code number reached: %s%d", prefix, maxNumber)
	}

	return fmt.Sprintf("%s%0*d", prefix, digits, nextNumber), nil
}
func GenerateCodeAllowance(prefix string, padding int, getLastCodeFn func() (string, error)) (string, error) {
	lastCode, err := getLastCodeFn()
	if err != nil {
		return "", fmt.Errorf("không thể lấy mã cuối: %w", err)
	}

	var newNumber int

	if lastCode == "" {
		newNumber = 1
	} else {
		if !strings.HasPrefix(lastCode, prefix) {
			return "", fmt.Errorf("invalid code format: %s, expected prefix: %s", lastCode, prefix)
		}

		numberPart := strings.TrimPrefix(lastCode, prefix)

		n, err := strconv.Atoi(numberPart)
		if err != nil {
			return "", fmt.Errorf("không thể parse số từ mã: %s", lastCode)
		}
		newNumber = n + 1
	}

	format := "%0" + strconv.Itoa(padding) + "d"
	newCode := fmt.Sprintf(prefix+format, newNumber)
	return newCode, nil
}
