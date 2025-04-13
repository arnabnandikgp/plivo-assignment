package db

import (
	"fmt"
	"log"
	"os"

	"github.com/status_page/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() {
	// --->>here<<--- Database connection is initialized using environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
}

// MigrateDB runs database migrations for all models
func MigrateDB() {
	// --->>here<<--- Database migrations are performed with GORM
	err := DB.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Service{},
		&models.Incident{},
		&models.IncidentUpdate{},
		&models.IncidentService{},
	)

	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	log.Println("Database migrations completed")
}
