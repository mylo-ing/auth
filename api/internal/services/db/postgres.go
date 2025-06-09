package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	var db *gorm.DB
	var err error

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSL_MODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)
	log.Printf("Connecting to DB host=%s user=%s dbname=%s port=%s sslmode=%s",
		host, user, dbname, port, sslmode)

	for i := 1; i <= 20; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("db not ready (try %d/20): %v", i, err)
			time.Sleep(3 * time.Second)
		}
	}

	return db
}
