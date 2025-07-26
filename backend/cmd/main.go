package main

import (
	_ "flow-sight-backend/docs"
	"flow-sight-backend/internal/api"
	"flow-sight-backend/internal/config"
	"flow-sight-backend/internal/database"
	"log"

	"github.com/joho/godotenv"
)

// @title Flow Sight API
// @version 1.0
// @description API for personal financial management application
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(cfg.Database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize and start API server
	server := api.NewServer(db, cfg)
	log.Printf("Starting server on port %s", cfg.Port)
	if err := server.Start(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
