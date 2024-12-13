package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aodihis/go-rest-signup-login/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDb() error {
	dbHost := config.GetEnv("DB_HOST")
	dbPort := config.GetEnv("DB_PORT")
	dbUser := config.GetEnv("DB_USER")
	dbPassword := config.GetEnv("DB_PASSWORD")
	dbName := config.GetEnv("DB_NAME")
	dbSSLMode := config.GetEnv("DB_SSLMODE")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	return nil
}

func CloseDb() {
	if DB != nil {
		DB.Close()
	}
}
