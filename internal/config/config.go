package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	StorageFile  string
	Environment  string
	LogLevel     string
	DATABASE_URL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	log.Println("DEBUG start")

	if err != nil {
		log.Println(".env tidak ditemukan, pakai environment sistem")
	}

	config := &Config{
		Port:         getEnv("PORT", "8080"),
		StorageFile:  getEnv("STORAGE_FILE", "storage.json"),
		Environment:  getEnv("ENV", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		DATABASE_URL: getEnv("DATABASE_URL", "postgres://postgres:admin@localhost:5433/task_api"),
	}

	if config.Port == "" {
		log.Fatal("PORT tidak boleh kosong")
	}

	return config
}

func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
