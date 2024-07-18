package config

import (
	"em-test/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	err = godotenv.Load()
	if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
  }

	host := os.Getenv("DB_HOST")
  user := os.Getenv("DB_USER")
  password := os.Getenv("DB_PASSWORD")
  dbname := os.Getenv("DB_NAME")
  port := os.Getenv("DB_PORT")
  sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        host, user, password, dbname, port, sslmode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
    log.Fatalf("failed to connect to database: %v", err)
  }

	fmt.Println("Database connected successfully")

	err = DB.AutoMigrate(&models.User{}, &models.Task{}, &models.TaskLog{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

  fmt.Println("Database migrated successfully")
}