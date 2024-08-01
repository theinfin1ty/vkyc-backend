package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {
	_ = godotenv.Load(".env")

	return os.Getenv(key)
}
