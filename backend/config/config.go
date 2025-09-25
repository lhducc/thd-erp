package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres Postgres
	CORS     CORS
}

// Server config struct
type ServerConfig struct {
	AppVersion string
	Port       string
	SSL        bool
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	TimeZone string
	DBSource string
}

type CORS struct {
	AllowOrigins []string
}

var AppConfig Config

func LoadConfig() {
	LoadEnv()

	// Load SSL
	if ssl := os.Getenv("SSL"); ssl != "" {
		parsed, err := strconv.ParseBool(ssl)
		if err != nil {
			log.Printf("Không thể parse SSL=%s thành bool, dùng mặc định false\n", ssl)
		}
		AppConfig.Server.SSL = parsed
		log.Printf("SSL loaded from .env: %v", AppConfig.Server.SSL)
	}

	// Load DB_SOURCE từ biến môi trường
	if dbSource := os.Getenv("DB_SOURCE"); dbSource != "" {
		AppConfig.Postgres.DBSource = dbSource
		log.Printf("DB_SOURCE loaded from .env: %s", AppConfig.Postgres.DBSource)
	}

	// Override CORS origins from environment variable if set
	if corsOrigins := os.Getenv("CORS_ALLOW_ORIGINS"); corsOrigins != "" {
		// Split comma-separated values and trim spaces
		origins := strings.Split(corsOrigins, ",")
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
		AppConfig.CORS.AllowOrigins = origins
		log.Printf("CORS origins loaded from environment: %v", origins)
	} else {
		log.Printf("CORS origins loaded from config file: %v", AppConfig.CORS.AllowOrigins)
	}
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy file .env hoặc lỗi khi load.")
	} else {
		log.Println(".env file loaded")
	}
}
