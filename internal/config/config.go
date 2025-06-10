package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	SMTPFrom    string
	SMTPPass    string
	SMTPHost    string
	SMTPPort    string
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Формируем строку подключения по умолчанию
		dbUser := "postgres"
		dbPassword := "admin"
		dbName := "findbgo1"
		dbHost := "localhost"
		dbPort := "5432"
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretkey123" // Значение по умолчанию для разработки
	}

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "noreply@financeapp.local"
	}
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		smtpPass = "testpass"
	}
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost"
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "1025"
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		SMTPFrom:    smtpFrom,
		SMTPPass:    smtpPass,
		SMTPHost:    smtpHost,
		SMTPPort:    smtpPort,
	}
}
