package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/sadamoto/kindle-notifier/internal/db"
	"github.com/sadamoto/kindle-notifier/internal/db/seeder"
)

func main() {
	// Parse command line flags
	clear := flag.Bool("clear", false, "Clear all seeded data")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database connection
	db.Initialize()

	// Create seeder registry
	registry := seeder.NewRegistry()

	// Register seeders (order matters for dependencies)
	registry.Register(seeder.NewCategorySeeder())
	registry.Register(seeder.NewBookSeeder())

	// Run seeders
	if *clear {
		if err := registry.ClearAll(); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		log.Println("Successfully cleared all seeded data")
	} else {
		if err := registry.SeedAll(); err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
		log.Println("Successfully seeded all data")
	}
} 