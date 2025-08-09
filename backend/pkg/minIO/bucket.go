package minIO

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"erp/backend/config"

	"github.com/minio/minio-go/v7"
)

type Bucket string

const (
	AttendanceBucket Bucket = "attendance"
)

func UploadImageToMinIO(
	ctx context.Context,
	bucketName Bucket,
	objectName string,
	file multipart.File,
	fileSize int64,
	contentType string,
	lifeCycleTimeDay int,
) error {
	minioClient := config.MinioClient
	bucketNameStr := string(bucketName)

	// Check bucket existence
	exists, err := minioClient.BucketExists(ctx, bucketNameStr)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketNameStr, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Println("Created bucket:", bucketName)

		if lifeCycleTimeDay > 0 {
			err = SetupLifecycle(minioClient, bucketName, lifeCycleTimeDay)
			if err != nil {
				log.Fatalf("Failed to setup lifecycle: %v", err)
			}
		} else {
			return fmt.Errorf("lifecycle time must be greater than 0")
		}
	}

	// Upload file
	_, err = minioClient.PutObject(ctx, bucketNameStr, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	return nil
}

func GeneratePresignedURL(
	ctx context.Context,
	bucketName Bucket,
	objectName string,
	expireTime time.Duration,
) (string, error) {
	minioClient := config.MinioClient

	// MinIO only allows max 7 days for presigned URLs
	if expireTime > 7*24*time.Hour {
		expireTime = 7 * 24 * time.Hour
	}

	bucketNameStr := string(bucketName)

	url, err := minioClient.PresignedGetObject(ctx, bucketNameStr, objectName, expireTime, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned URL failed: %w", err)
	}

	if os.Getenv("APP_ENV") == "docker" {
		finalURL := strings.Replace(url.String(), os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_PUBLIC_ENDPOINT"), 1)
		return finalURL, nil
	}
	return url.String(), nil
}

// auto delete file in bucket after lifeTimeDay
func SetupLifecycle(minioClient *minio.Client, bucketName Bucket, lifeTimeDay int) error {
	ctx := context.Background()

	bucketNameStr := string(bucketName)

	cfg := lifecycle.NewConfiguration()
	ruleID := fmt.Sprintf("expire-after-%d-days", lifeTimeDay)

	cfg.Rules = []lifecycle.Rule{
		{
			ID:     ruleID,
			Status: "Enabled",
			Prefix: "",
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(lifeTimeDay),
			},
		},
	}

	err := minioClient.SetBucketLifecycle(ctx, bucketNameStr, cfg)
	if err != nil {
		return err
	}

	log.Printf("Đã thiết lập lifecycle rule cho bucket %s", bucketName)
	return nil
}
