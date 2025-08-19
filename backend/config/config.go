package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Lỗi lấy working directory: %v", err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(wd, "config"))
	viper.AddConfigPath(filepath.Join(wd, "..", "..", "config"))

	// Đọc file config, nếu lỗi thì dừng chương trình
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Lỗi đọc config file: %v", err)
	}

	// Parse dữ liệu từ file config vào biến AppConfig
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
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

//func LoadEnv() {
//	envPath := filepath.Join("..", ".env")
//
//	if err := godotenv.Load(envPath); err != nil {
//		log.Printf("Không tìm thấy file .env tại %s: %v", envPath, err)
//	} else {
//		log.Printf("Đã load file .env từ: %s", envPath)
//	}
//}

//func LoadConfig() {
//	// Load biến môi trường từ file .env ở thư mục cha
//	LoadEnv()
//
//	// Load config từ file YML (nếu cần)
//	wd, err := os.Getwd()
//	if err != nil {
//		log.Fatalf("Lỗi lấy working directory: %v", err)
//	}
//
//	viper.SetConfigName("config")
//	viper.SetConfigType("yml")
//	viper.AddConfigPath(filepath.Join(wd, "config"))
//	viper.AddConfigPath(filepath.Join(wd, "..", "..", "config"))
//
//	if err := viper.ReadInConfig(); err != nil {
//		log.Printf("Không đọc được file config.yml: %v", err)
//	} else {
//		if err := viper.Unmarshal(&AppConfig); err != nil {
//			log.Printf("Lỗi parse config.yml: %v", err)
//		}
//	}
//
//	// Luôn tạo DSN từ biến môi trường
//	AppConfig.Postgres.DBSource = BuildDSN()
//	log.Printf("DSN kết nối PostgreSQL: %s", maskPassword(AppConfig.Postgres.DBSource))
//}

//// Hàm helper che giấu password trong log
//func maskPassword(dsn string) string {
//	if strings.Contains(dsn, "password=") {
//		return regexp.MustCompile(`password=[^ ]+`).ReplaceAllString(dsn, "password=*****")
//	}
//	return dsn
//}

//func getEnvValue(envKey, defaultValue string) string {
//	if value := os.Getenv(envKey); value != "" {
//		return value
//	}
//	return defaultValue
//}

//func BuildDSN() string {
//	// Lấy giá trị từ biến môi trường (ưu tiên) hoặc config.yml
//	host := getEnvValue("DB_HOST", AppConfig.Postgres.Host)
//	port := getEnvValue("DB_PORT", AppConfig.Postgres.Port)
//	user := getEnvValue("DB_USER", AppConfig.Postgres.User)
//	password := getEnvValue("DB_PASSWORD", AppConfig.Postgres.Password)
//	dbname := getEnvValue("DB_NAME", AppConfig.Postgres.DBName)
//	timezone := getEnvValue("DB_TIMEZONE", "Asia/Ho_Chi_Minh")
//
//	// Xây dựng chuỗi DSN
//	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
//		host,
//		port,
//		user,
//		password,
//		dbname,
//		timezone)
//}
