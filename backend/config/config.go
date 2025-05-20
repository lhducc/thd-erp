package config

import (
	"log"
	"os"
	"path/filepath"

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
}
