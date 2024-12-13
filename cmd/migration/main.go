package main

import (
	"log"

	"github.com/aodihis/go-rest-signup-login/config"
	"github.com/aodihis/go-rest-signup-login/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadEnv()
	database.ConnectDb()
	defer database.CloseDb()

	driver, err := postgres.WithInstance(database.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create db driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path to the migrations folder
		"postgres",          // Database name
		driver,
	)

	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	log.Println("Applying migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied!")
}
