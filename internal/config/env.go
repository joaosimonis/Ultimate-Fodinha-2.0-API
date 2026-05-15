package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AllowedOrigins      string
	JWTSecret           string
	JWTExpiresInSeconds int
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	dbHost, err := requiredEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	dbUser, err := requiredEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}

	dbPassword, err := requiredEnv("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	dbName, err := requiredEnv("DB_NAME")
	if err != nil {
		return Config{}, err
	}

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	jwtExpiresInSeconds, err := parseIntEnv("JWT_EXPIRES_IN_SECONDS", 3600)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppPort:             getEnv("APP_PORT", getEnv("PORT", "8080")),
		AllowedOrigins:      getEnv("ALLOWED_ORIGINS", ""),
		JWTSecret:           jwtSecret,
		JWTExpiresInSeconds: jwtExpiresInSeconds,
		DBHost:              dbHost,
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              dbUser,
		DBPassword:          dbPassword,
		DBName:              dbName,
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	return value
}

func requiredEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}

	return value, nil
}

func parseIntEnv(key string, fallback int) (int, error) {
	rawValue := getEnv(key, "")
	if rawValue == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("environment variable %s must be a positive integer", key)
	}

	return value, nil
}
