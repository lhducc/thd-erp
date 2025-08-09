package config

import (
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

type MinIOConfig struct {
	PublicMinioEndPoint string
	MinioEndpoint       string
	AccessKey           string
	SecretKey           string
	UseSSL              bool
}

var MinIO MinIOConfig

var MinioClient *minio.Client

//var MinioPublicURLImage *minio.Client

func LoadMinIOConfig() {
	if os.Getenv("APP_ENV") != "docker" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Lỗi load .env")
		}
	}
	MinIO.MinioEndpoint = os.Getenv("MINIO_ENDPOINT")
	MinIO.PublicMinioEndPoint = os.Getenv("MINIO_PUBLIC_ENDPOINT")
	MinIO.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	MinIO.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	MinIO.UseSSL = false

	internalClient, err := minio.New(MinIO.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinIO.AccessKey, MinIO.SecretKey, ""),
		Secure: MinIO.UseSSL,
	})
	if err != nil {
		log.Fatalf("Lỗi tạo MinIO client: %v", err)
	}

	MinioClient = internalClient
}
