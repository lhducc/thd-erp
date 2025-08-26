package utils

import (
	"fmt"
	"time"
)

func GenerateFileName(employeeID string, timestamp time.Time, filename string) string {
	return fmt.Sprintf("%s_%d_%s", employeeID, timestamp.Unix(), filename)
}
