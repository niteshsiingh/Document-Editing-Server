package config

import (
	"errors"
	"os"
)

type Env struct {
	NodeEnv                 string
	Host                    string
	Port                    string
	DatabaseURL             string
	User                    string
	Password                string
	DBHost                  string
	DBPort                  string
	Database                string
	SMTPHost                string
	SMTPUser                string
	SMTPPassword            string
	AccessTokenSecret       string
	AccessTokenExpiration   string
	RefreshTokenSecret      string
	RefreshTokenExpiration  string
	VerifyEmailSecret       string
	PasswordResetSecret     string
	PasswordResetExpiration string
	FrontEndURL             string
}

func InitEnv() (*Env, error) {
	nodeEnv := os.Getenv("NODE_ENV")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DATABASE_URL")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	database := os.Getenv("DATABASE")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	accessTokenExpiration := os.Getenv("ACCESS_TOKEN_EXPIRATION")
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	refreshTokenExpiration := os.Getenv("REFRESH_TOKEN_EXPIRATION")
	verifyEmailSecret := os.Getenv("VERIFY_EMAIL_SECRET")
	passwordResetSecret := os.Getenv("PASSWORD_RESET_SECRET")
	passwordResetExpiration := os.Getenv("PASSWORD_RESET_EXPIRATION")
	frontEndURL := os.Getenv("FRONT_END_URL")

	if nodeEnv == "" ||
		host == "" ||
		port == "" ||
		databaseURL == "" ||
		user == "" ||
		dbHost == "" ||
		dbPort == "" ||
		database == "" ||
		smtpHost == "" ||
		smtpUser == "" ||
		smtpPassword == "" ||
		accessTokenSecret == "" ||
		accessTokenExpiration == "" ||
		refreshTokenSecret == "" ||
		refreshTokenExpiration == "" ||
		verifyEmailSecret == "" ||
		passwordResetSecret == "" ||
		passwordResetExpiration == "" ||
		frontEndURL == "" {
		return nil, errors.New("environment variables missing")
	}

	env := &Env{
		NodeEnv:                 nodeEnv,
		Host:                    host,
		Port:                    port,
		DatabaseURL:             databaseURL,
		User:                    user,
		Password:                password,
		DBHost:                  dbHost,
		DBPort:                  dbPort,
		Database:                database,
		SMTPHost:                smtpHost,
		SMTPUser:                smtpUser,
		SMTPPassword:            smtpPassword,
		AccessTokenSecret:       accessTokenSecret,
		AccessTokenExpiration:   accessTokenExpiration,
		RefreshTokenSecret:      refreshTokenSecret,
		RefreshTokenExpiration:  refreshTokenExpiration,
		VerifyEmailSecret:       verifyEmailSecret,
		PasswordResetSecret:     passwordResetSecret,
		PasswordResetExpiration: passwordResetExpiration,
		FrontEndURL:             frontEndURL,
	}

	return env, nil
}
