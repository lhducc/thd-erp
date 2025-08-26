package minIO

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/minio/minio-go/v7/pkg/lifecycle"

	"erp/backend/config"

	"github.com/minio/minio-go/v7"
)

type Bucket string

const (
	AttendanceBucket Bucket = "attendance"
	DecisionBucket   Bucket = "decision"
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

		// Set bucket policy to allow public read access
		err = SetBucketPublicReadPolicy(minioClient, bucketName)
		if err != nil {
			log.Printf("Warning: Failed to set public read policy for bucket %s: %v", bucketName, err)
		}

		if lifeCycleTimeDay > 0 {
			err = SetupLifecycle(minioClient, bucketName, lifeCycleTimeDay)
			if err != nil {
				log.Fatalf("Failed to setup lifecycle: %v", err)
			}
		} else {
			return fmt.Errorf("lifecycle time must be greater than 0")
		}
	} else {
		// Bucket exists, ensure it has public read policy
		err = SetBucketPublicReadPolicy(minioClient, bucketName)
		if err != nil {
			log.Printf("Warning: Failed to set public read policy for existing bucket %s: %v", bucketName, err)
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
	// Try presigned URL first using URL client (if available)
	urlClient := config.MinioURLClient
	bucketNameStr := string(bucketName)

	// MinIO only allows max 7 days for presigned URLs
	if expireTime > 7*24*time.Hour {
		expireTime = 7 * 24 * time.Hour
	}

	if urlClient != nil {
		if u, err := urlClient.PresignedGetObject(ctx, bucketNameStr, objectName, expireTime, nil); err == nil {
			return u.String(), nil
		} else {
			log.Printf("Warning: presigned URL generation failed for %s/%s: %v", bucketNameStr, objectName, err)
		}
	} else {
		log.Printf("MinIO URL client not available, falling back to direct object URL for %s/%s", bucketNameStr, objectName)
	}

	// Fallback: construct a direct object URL from configured public endpoint or internal endpoint
	endpoint := strings.TrimRight(config.MinIO.PublicMinioEndPoint, "/")
	if endpoint == "" {
		endpoint = strings.TrimRight(config.MinIO.MinioEndpoint, "/")
	}
	if endpoint == "" {
		return "", fmt.Errorf("no MinIO endpoint configured for fallback URL")
	}

	// Ensure scheme present
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		if config.MinIO.UseSSL {
			endpoint = "https://" + endpoint
		} else {
			endpoint = "http://" + endpoint
		}
	}

	return fmt.Sprintf("%s/%s/%s", endpoint, bucketNameStr, objectName), nil
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

// SetBucketPublicReadPolicy sets the bucket policy to allow public read access
func SetBucketPublicReadPolicy(minioClient *minio.Client, bucketName Bucket) error {
	ctx := context.Background()
	bucketNameStr := string(bucketName)

	// Use MinIO's predefined public read policy
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + bucketNameStr + `/*"]
			}
		]
	}`

	err := minioClient.SetBucketPolicy(ctx, bucketNameStr, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	log.Printf("Set public read policy for bucket %s", bucketName)
	return nil
}

// EnsureBucketPublicAccess ensures that the bucket has public read access
func EnsureBucketPublicAccess(bucketName Bucket) error {
	minioClient := config.MinioClient
	return SetBucketPublicReadPolicy(minioClient, bucketName)
}
