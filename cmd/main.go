package main

import (
	"log"

	"github.com/sadamoto/kindle-notifier/internal/api"
	"github.com/sadamoto/kindle-notifier/internal/db"
	"github.com/sadamoto/kindle-notifier/internal/job"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database connection
	db.Initialize()

	// Create PA-API client
	client := api.NewPAAPIClient()

	// Create and run the book checker job
	checker := job.NewBookChecker(client)
	if err := checker.Run(); err != nil {
		log.Printf("Error running book checker job: %v", err)
	}
} 