package config

import (
    "fmt"
    "os"
    "retail_pulse_project/models"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

    dsn := os.Getenv("DATABASE_URL")
    fmt.Println("Connecting to database with DSN:", dsn)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
    if err != nil {
        fmt.Println("Failed to connect to the database:", err)
        panic("Cannot proceed without database connection!")
    }

    // Automatically create tables based on models
    err = db.AutoMigrate(&models.Job{}, &models.Image{}, &models.Store{})
    if err != nil {
        fmt.Println("Failed to run migrations:", err)
        panic("Cannot proceed without running migrations!")
    }

    DB = db
    fmt.Println("Database connected and tables migrated successfully!")
}
