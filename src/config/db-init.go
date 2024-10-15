package config

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niteshsiingh/doc-server/src/database/tables/databases"
)

func InitializeDatabase() (*databases.Queries, error) {
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

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	// defer conn.Close(ctx)
	db := databases.New(conn)
	log.Println("Database connection established")
	return db, nil
}
