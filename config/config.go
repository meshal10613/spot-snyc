package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all environment-level configuration.
type Config struct {
	Port          string
	DSN           string
	JWTSecret     string
	AdminName     string
	AdminEmail    string
	AdminPassword string
}

// LoadEnv reads the .env file and returns a Config struct.
func LoadEnv() (*Config, error) {
	// In production .env may not exist — ignore the error
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		adminName = "Admin"
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Fatal("ADMIN_EMAIL environment variable is required")
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is required")
	}

	return &Config{
		Port:          port,
		DSN:           dsn,
		JWTSecret:     jwtSecret,
		AdminName:     adminName,
		AdminEmail:    adminEmail,
		AdminPassword: adminPassword,
	}, nil
}
