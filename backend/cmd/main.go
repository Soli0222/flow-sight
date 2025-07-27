package main

import (
	"context"

	_ "github.com/Soli0222/flow-sight/backend/docs"
	"github.com/Soli0222/flow-sight/backend/internal/api"
	"github.com/Soli0222/flow-sight/backend/internal/config"
	"github.com/Soli0222/flow-sight/backend/internal/database"
	"github.com/Soli0222/flow-sight/backend/internal/logger"
	"github.com/Soli0222/flow-sight/backend/internal/version"

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
		// Note: .env file is optional, so we don't log this as an error
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.New(cfg, version.Version)
	ctx := context.Background()

	appLogger.InfoContext(ctx, "Starting Flow Sight Backend",
		"version", version.Version,
		"environment", cfg.Env,
	)

	// Initialize database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		appLogger.ErrorContext(ctx, "Failed to connect to database", "error", err.Error())
		return
	}
	defer db.Close()

	appLogger.InfoContext(ctx, "Database connected successfully")

	// Run migrations
	if err := database.Migrate(cfg.Database); err != nil {
		appLogger.ErrorContext(ctx, "Failed to run migrations", "error", err.Error())
		return
	}

	appLogger.InfoContext(ctx, "Database migrations completed")

	// Initialize and start API server
	server := api.NewServer(db, cfg, appLogger)
	appLogger.InfoContext(ctx, "Starting server", "port", cfg.Port)

	if err := server.Start(":" + cfg.Port); err != nil {
		appLogger.ErrorContext(ctx, "Failed to start server", "error", err.Error())
		return
	}
}
