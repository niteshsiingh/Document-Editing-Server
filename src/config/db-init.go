package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDatabase() (*gorm.DB, error) {
	env := os.Getenv("ENV")
	dbURL := os.Getenv("DB_URL")

	var dsn string
	if env == "prod" {
		dsn = fmt.Sprintf("%s?sslmode=require&sslrootcert=/path/to/ca.crt", dbURL)
	} else {
		dbUser := os.Getenv("USER")
		dbPass := os.Getenv("PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DATABASE")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connection established")
	return db, nil
}
