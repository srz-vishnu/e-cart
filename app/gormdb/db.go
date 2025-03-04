package gormdb

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() (*gorm.DB, error) {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	fmt.Println("Connection string:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connection error", err)
	}

	err1 := Automigration(db)
	if err1 != nil {
		log.Fatalf("error on automigration %v", err1)
	}
	log.Println("Migration successfully completed")

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal("connection error", err)
	}

	// Test the connection
	err = sqlDb.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}

	fmt.Println("Successfully connected to the database!!!")

	return db, nil
}
