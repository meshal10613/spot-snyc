package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all environment-level configuration.
type Config struct {
	Port      string
	DSN       string
	JWTSecret string
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

	return &Config{
		Port:      port,
		DSN:       dsn,
		JWTSecret: jwtSecret,
	}, nil
}
