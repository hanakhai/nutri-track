package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func LoadConfig() *Config {
	// Membaca file .env jika ada (tidak error jika file tidak ada di environment tertentu)
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca konfigurasi dari system environment.")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Fatal: DATABASE_URL wajib diisi pada environment!")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Fatal: JWT_SECRET wajib diisi pada environment!")
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	return &Config{
		AppPort:     appPort,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		Environment: environment,
	}
}

