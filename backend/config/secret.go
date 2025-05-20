package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const SecretKeyEnv = "SECRET_KEY"

func GetSecretKey() string {
	_ = godotenv.Load()

	key := os.Getenv(SecretKeyEnv)
	if key == "" {
		key = generateSecretKey(32) // 32 bytes = 64 chars hex
		err := saveKeyToEnv(key)
		if err != nil {
			log.Fatalf("Failed to save secret key to .env: %v", err)
		}
	}
	return key
}

func generateSecretKey(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate secret key: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func saveKeyToEnv(key string) error {
	f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(SecretKeyEnv + "=" + key + "\n")
	return err
}

func updateEnvFile(key, value string) {
	filePath := ".env"

	data, err := os.ReadFile(filePath)
	if err != nil {
		data = []byte{}
	}

	lines := strings.Split(string(data), "\n")
	found := false

	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, key+"="+value)
	}

	newData := strings.Join(lines, "\n")

	err = os.WriteFile(filePath, []byte(newData), 0644)
	if err != nil {
		log.Fatalf("Failed to write .env file: %v", err)
	}
}
