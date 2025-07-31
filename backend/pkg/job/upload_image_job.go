// internal/pkg/job/upload_image_job.go
package job

import (
	"mime/multipart"
	"time"
)

type UploadJob struct {
	EmployeeID string
	Timestamp  time.Time
	File       multipart.File
	Header     *multipart.FileHeader
	FileName   string
}
