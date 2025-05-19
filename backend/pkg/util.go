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

func GenerateCode(prefix string, digits int, getLastCodeFunc func() (string, error)) (string, error) {
	lastCode, err := getLastCodeFunc()
	if err != nil {
		if err.Error() == "record not found" {
			return fmt.Sprintf("%s%0*d", prefix, digits, 1), nil
		}
		return "", fmt.Errorf("failed to get last code: %w", err)
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
