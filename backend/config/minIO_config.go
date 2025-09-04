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

	// Determine whether to use SSL. Prefer explicit env var MINIO_USE_SSL if provided.
	// Otherwise default to false when running in docker or when endpoint looks like internal docker name.
	sslEnv := os.Getenv("MINIO_USE_SSL")
	if sslEnv != "" {
		MinIO.UseSSL = strings.ToLower(sslEnv) == "true"
	} else {
		// If running in docker or endpoint points to internal service, assume HTTP (no SSL) by default.
		if os.Getenv("APP_ENV") == "docker" || strings.Contains(MinIO.MinioEndpoint, "minio:") || strings.Contains(MinIO.MinioEndpoint, "192.168.1.58:9000") {
			MinIO.UseSSL = false
		} else {
			MinIO.UseSSL = true
		}
	}

	// Create transport. When using SSL skip verification only if SSL is enabled.
	var tr *http.Transport
	if MinIO.UseSSL {
		tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	} else {
		tr = &http.Transport{}
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
	// URL transport: preserve the custom dialer redirection, and set TLS config only when SSL is enabled.
	urlTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Allow self-signed certificates
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Redirect localhost:9000 to the actual MinIO container via host IP
			if strings.Contains(addr, "192.168.1.58:9000") {
				addr = "192.168.1.58:9000"
			}
			return customDialer.DialContext(ctx, network, addr)
		},
	}
	if MinIO.UseSSL {
		urlTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	urlClient, err := minio.New("192.168.1.58:9000", &minio.Options{
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
