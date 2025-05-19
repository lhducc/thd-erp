package config

import (
	"log"

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
	DBSource string
}

type CORS struct {
	AllowOrigins []string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("C:\\thd\\erp\\erp\\backend\\config")

	// Đọc file config, nếu lỗi thì dừng chương trình
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Lỗi đọc config file: %v", err)
	}

	// Parse dữ liệu từ file config vào biến AppConfig
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
