package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/defi-dashboard/backend/internal/config"
	"github.com/defi-dashboard/backend/internal/jobs"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/external"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)
	logger.Info("Starting DeFi Dashboard Worker", "version", "1.0.0")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Database connection
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer dbpool.Close()

	// Test database connection
	if err := dbpool.Ping(ctx); err != nil {
		logger.Fatal("Failed to ping database", "error", err)
	}
	logger.Info("Successfully connected to database")

	// Initialize external API clients
	coinGeckoClient := external.NewCoinGeckoClient(cfg.CoinGeckoAPIKey)
	defiLlamaClient := external.NewDefiLlamaClient()

	// Initialize repositories
	alertRepo := repos.NewAlertRepository(dbpool)
	userRepo := repos.NewUserRepository(dbpool)

	// Initialize services
	alertService := services.NewAlertService(alertRepo, userRepo)

	// Initialize job handlers
	priceJob := jobs.NewPriceRefreshJob(dbpool, coinGeckoClient, defiLlamaClient)
	alertJob := jobs.NewAlertEvaluatorJob(dbpool, alertService, alertRepo)

	// Create cron scheduler with seconds support
	c := cron.New(cron.WithSeconds())

	// Schedule jobs
	// Price refresh every 10 minutes
	_, err = c.AddFunc("0 */10 * * * *", func() {
		runJob(ctx, "price-refresh", priceJob.Run)
	})
	if err != nil {
		logger.Fatal("Failed to schedule price refresh job", "error", err)
	}

	// Alert evaluation every 5 minutes
	_, err = c.AddFunc("0 */5 * * * *", func() {
		runJob(ctx, "alert-evaluator", alertJob.Run)
	})
	if err != nil {
		logger.Fatal("Failed to schedule alert evaluator job", "error", err)
	}

	// Start cron scheduler
	c.Start()
	logger.Info("Worker scheduled jobs started")

	// Run initial jobs on startup
	logger.Info("Running initial jobs on startup")
	runJob(ctx, "price-refresh-startup", priceJob.Run)
	runJob(ctx, "alert-evaluator-startup", alertJob.Run)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping worker...")

	// Stop cron scheduler
	cronCtx := c.Stop()
	
	// Cancel context to stop any running jobs
	cancel()

	// Wait for cron jobs to finish with timeout
	done := make(chan struct{})
	go func() {
		<-cronCtx.Done()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("All jobs stopped gracefully")
	case <-time.After(30 * time.Second):
		logger.Warn("Timeout waiting for jobs to stop")
	}

	logger.Info("Worker shutdown complete")
}

// runJob executes a job with proper error handling and logging
func runJob(ctx context.Context, jobName string, jobFunc func(context.Context) error) {
	logger.Info("Starting job", "job", jobName)
	start := time.Now()

	// Create a timeout context for the job
	jobCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Run the job
	if err := jobFunc(jobCtx); err != nil {
		logger.Error("Job failed", 
			"job", jobName, 
			"error", err, 
			"duration", time.Since(start))
		return
	}

	logger.Info("Job completed successfully", 
		"job", jobName, 
		"duration", time.Since(start))
}