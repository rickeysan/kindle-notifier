package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
	mu   sync.RWMutex
)

// GetDB returns the database connection, initializing it if necessary
func GetDB() *gorm.DB {
	mu.RLock()
	if DB != nil {
		defer mu.RUnlock()
		return DB
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	
	if DB != nil {
		return DB
	}

	if err := Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		return nil
	}
	
	return DB
}

// Initialize initializes the database connection
func Initialize() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// IsConnected checks if the database connection is established
func IsConnected() bool {
	if DB == nil {
		return false
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}
	
	return sqlDB.Ping() == nil
} 