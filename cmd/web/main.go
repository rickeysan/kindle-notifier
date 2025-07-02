package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sadamoto/kindle-notifier/internal/web"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// Basic health check that doesn't require DB
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Load environment variables in development
	if os.Getenv("RENDER") == "" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	// Create handler
	handler, err := web.NewHandler()
	if err != nil {
		log.Printf("Warning: Error creating handler: %v", err)
		// Continue anyway, as we want the health check endpoint to work
	}

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheck)
	
	// Only set up the main handler if we could create it
	if handler != nil {
		mux.HandleFunc("/", handler.HandleIndex)
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Create custom server with timeouts
	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%s", port),
		Handler:           mux,
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server
	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
} 