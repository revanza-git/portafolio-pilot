package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/defi-dashboard/backend/internal/config"
	"github.com/defi-dashboard/backend/internal/router"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Database connection
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer dbpool.Close()

	// Test database connection
	if err := dbpool.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", "error", err)
	}

	logger.Info("Successfully connected to database")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "DeFi Dashboard API",
		ErrorHandler:          router.CustomErrorHandler,
		ReadTimeout:           time.Second * 30,
		WriteTimeout:          time.Second * 30,
		IdleTimeout:           time.Second * 30,
		DisableStartupMessage: true,
	})

	// Setup routes
	router.SetupRoutes(app, dbpool, cfg)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			logger.Error("Server shutdown error", "error", err)
		}
	}()

	// Start server
	logger.Info("Starting server", "port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}