package job

import (
	"context"
	"erp/backend/pkg/minIO"
	"log"
	"mime/multipart"
	"time"
)

var JobQueue chan UploadJob
var FailedJobQueue chan UploadJob

func InitWorkerPool(poolSize int) {
	JobQueue = make(chan UploadJob, 1000)
	FailedJobQueue = make(chan UploadJob, 100)

	for i := 0; i < poolSize; i++ {
		go worker(i)
	}

	go handleFailedJobs()
}

func worker(id int) {
	for job := range JobQueue {
		var err error
		maxRetries := 3

		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = handleUpload(job)
			if err == nil {
				log.Printf("[Worker %d] Upload success: %s\n", id, job.Header.Filename)
				break
			}

			log.Printf("[Worker %d] Upload failed (attempt %d/%d): %v\n", id, attempt, maxRetries, err)
			time.Sleep(500 * time.Millisecond)
		}

		if err != nil {
			log.Printf("[Worker %d] Upload permanently failed: %s\n", id, job.Header.Filename)
			FailedJobQueue <- job
		}
	}
}

// Hàm xử lý upload thực tế
func handleUpload(job UploadJob) error {
	defer func() {
		if cerr := job.File.Close(); cerr != nil {
			log.Printf("⚠️ Error closing file: %v", cerr)
		}
	}()

	return minIO.UploadImageToMinIO(
		context.Background(),
		minIO.AttendanceBucket,
		job.FileName,
		job.File,
		job.Header.Size,
		job.Header.Header.Get("Content-Type"),
		365,
	)
}

// Đưa job vào hàng đợi chính
func EnqueueUploadJob(file multipart.File, employee string, header *multipart.FileHeader, timestamp time.Time, imageName string) {
	job := UploadJob{
		EmployeeID: employee,
		Timestamp:  timestamp,
		File:       file,
		Header:     header,
		FileName:   imageName,
	}
	JobQueue <- job
}

// Goroutine xử lý job lỗi
func handleFailedJobs() {
	for job := range FailedJobQueue {
		log.Printf("Failed job captured: %s by %s at %s", job.Header.Filename, job.EmployeeID, job.Timestamp.Format(time.RFC3339))
		// TODO: Ghi vào DB, gửi email cảnh báo, lưu log file,...
	}
}
