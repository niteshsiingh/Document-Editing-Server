package config

import (
	"errors"
	"os"
	"strconv"
)

type SMTP struct {
	Port         int
	SMTPHost     string
	SMTPUser     string
	SMTPPassword string
	Secure       bool
}

func InitSMTP() (*SMTP, error) {
	portStr := "587"
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	secureStr := os.Getenv("SMTP_SECURE")
	var secure bool
	if portStr == "" ||
		smtpHost == "" ||
		smtpUser == "" ||
		smtpPassword == "" ||
		secureStr == "" {
		return nil, errors.New("environment variables missing")
	}
	if secureStr == "true" {
		secure = true
	} else {
		secure = false
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	smtp := &SMTP{
		Port:         port,
		SMTPHost:     smtpHost,
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
		Secure:       secure,
	}

	return smtp, nil
}
