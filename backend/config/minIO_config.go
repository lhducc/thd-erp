package config

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
var MinioURLClient *minio.Client // Separate client for URL generation

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
	MinIO.UseSSL = true // Enable SSL for HTTPS MinIO

	// Create transport to skip SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Internal client for operations (uses Docker service name)
	internalClient, err := minio.New(MinIO.MinioEndpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(MinIO.AccessKey, MinIO.SecretKey, ""),
		Secure:    MinIO.UseSSL,
		Transport: tr,
	})
	if err != nil {
		log.Fatalf("Lỗi tạo MinIO internal client: %v", err)
	}

	// URL client for generating presigned URLs (uses localhost but redirects to host IP)
	customDialer := &net.Dialer{}
	urlTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Redirect localhost:9000 to the actual MinIO container via host IP
			if strings.Contains(addr, "localhost:9000") {
				addr = "172.18.0.1:9000"
			}
			return customDialer.DialContext(ctx, network, addr)
		},
	}

	urlClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:     credentials.NewStaticV4(MinIO.AccessKey, MinIO.SecretKey, ""),
		Secure:    MinIO.UseSSL,
		Transport: urlTransport,
	})
	if err != nil {
		log.Printf("Lỗi tạo MinIO URL client: %v, sử dụng internal client", err)
		MinioURLClient = internalClient
	} else {
		MinioURLClient = urlClient
	}

	log.Printf("Kết nối minio thành công")
	MinioClient = internalClient
}
