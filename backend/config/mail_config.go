package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var Mail MailConfig

func LoadMailConfig() {
	_ = godotenv.Load(".env")

	Mail.Host = os.Getenv("SMTP_HOST")
	Mail.Username = os.Getenv("SMTP_USERNAME")
	Mail.Password = os.Getenv("SMTP_PASSWORD")
	Mail.From = os.Getenv("EMAIL_FROM")

	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %s", portStr)
	}
	Mail.Port = port
}
