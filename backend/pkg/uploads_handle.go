package utils

import (
	"erp/backend/pkg/job"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func GenerateImageName(employeeID string, timestamp time.Time, filename string) string {
	return fmt.Sprintf("%s_%d_%s", employeeID, timestamp.Unix(), filename)
}

func HandleImageUpload(c *gin.Context, employeeID string, timestamp time.Time) (string, error) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		return "", err
	}

	// Tạo tên ảnh
	imageName := GenerateImageName(employeeID, timestamp, header.Filename)

	// Tạo job upload và đẩy vào queue
	jobData := job.UploadJob{
		EmployeeID: employeeID,
		Timestamp:  timestamp,
		File:       file,
		Header:     header,
	}
	job.JobQueue <- jobData

	return imageName, nil
}
