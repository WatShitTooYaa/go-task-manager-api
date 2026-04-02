package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	StorageFile string
	Environment string
	LogLevel    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	log.Println("DEBUG start")
	if err != nil {
		log.Println("tidak error saat load env")
		return nil
	}
	log.Fatal("file is nil")
	return &Config{
		Port:        getEnv("PORT", "8080"),
		StorageFile: getEnv("STORAGE_FILE", "storage.json"),
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
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
