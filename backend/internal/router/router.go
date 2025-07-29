package router

import (
	"time"

	"github.com/defi-dashboard/backend/internal/config"
	"github.com/defi-dashboard/backend/internal/handlers"
	"github.com/defi-dashboard/backend/internal/middleware"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/pnl"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CustomErrorHandler handles all errors in a consistent format
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var response interface{}

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		response = errors.New("FIBER_ERROR", e.Message, code)
	} else if e, ok := err.(*errors.AppError); ok {
		code = e.Status
		response = e
	} else {
		response = errors.Internal(err.Error())
	}

	logger.Error("Request error",
		"path", c.Path(),
		"method", c.Method(),
		"status", code,
		"error", err.Error(),
		"request_id", c.Locals("requestid"),
	)

	return c.Status(code).JSON(response)
}

func SetupRoutes(app *fiber.App, db *pgxpool.Pool, cfg *config.Config) {
	// Global middleware
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(recover.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Authorization,Content-Type,Accept,Origin,X-Requested-With,X-Alchemy-API-Key,X-CoinGecko-API-Key,X-Etherscan-API-Key,X-Infura-API-Key,x-alchemy-api-key,x-coingecko-api-key,x-etherscan-api-key,x-infura-api-key",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return errors.New("RATE_LIMIT_EXCEEDED", "Too many requests", fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	}))

	// Request logging middleware
	app.Use(middleware.RequestLogger())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		// TODO: Add database health check
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now().Unix(),
		})
	})

	// Initialize repositories
	userRepo := repos.NewUserRepository(db)
	walletRepo := repos.NewWalletRepository(db)
	tokenRepo := repos.NewTokenRepository(db)
	transactionRepo := repos.NewTransactionRepository(db)
	nonceRepo := repos.NewNonceRepository(db)
	
	// Yield repositories
	protocolRepo := repos.NewProtocolRepository(db)
	yieldPoolRepo := repos.NewYieldPoolRepository(db)
	yieldPositionRepo := repos.NewYieldPositionRepository(db)

	// Initialize services (blockchain services will be created dynamically with user API keys)
	authService := services.NewAuthService(userRepo, walletRepo, cfg.JWTSecret, cfg.JWTExpiry)
	siweService := services.NewSIWEService(userRepo, nonceRepo, "localhost") // TODO: Use actual domain from config
	portfolioService := services.NewPortfolioService(walletRepo, tokenRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	
	// Initialize bridge and swap services with external API clients
	bridgeService := services.NewBridgeService(
		cfg.GetLiFiClientConfig(),
		cfg.GetSocketClientConfig(),
	)
	swapService := services.NewSwapService(
		cfg.GetZeroXClientConfig(),
		cfg.GetOneInchClientConfig(),
	)
	
	yieldService := services.NewYieldService(yieldPoolRepo, yieldPositionRepo, protocolRepo, userRepo)
	
	// Initialize PnL service
	pnlRepo := pnl.NewRepository(db)
	pnlService := pnl.NewService(pnlRepo, walletRepo, tokenRepo)
	csvExporter := pnl.NewCSVExporter("/tmp") // TODO: Use configurable temp directory

	// Initialize Alert service
	alertRepo := repos.NewAlertRepository(db)
	alertService := services.NewAlertService(alertRepo, userRepo)

	// Initialize Watchlist repository
	watchlistRepo := repos.NewWatchlistRepository(db)

	// Initialize Admin repositories
	featureFlagRepo := repos.NewFeatureFlagRepository(db)
	systemBannerRepo := repos.NewSystemBannerRepository(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, siweService, cfg.JWTSecret, cfg.JWTExpiry)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	bridgeHandler := handlers.NewBridgeHandler(bridgeService)
	swapHandler := handlers.NewSwapHandler(swapService)
	yieldHandler := handlers.NewYieldHandler(yieldService)
	analyticsHandler := handlers.NewAnalyticsHandler(pnlService, csvExporter)
	alertHandler := handlers.NewAlertHandler(alertService)
	watchlistHandler := handlers.NewWatchlistHandler(watchlistRepo)
	adminHandler := handlers.NewAdminHandler(userRepo, featureFlagRepo, systemBannerRepo)

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes (no auth required)
	auth := v1.Group("/auth")
	
	// SIWE Authentication
	siwe := auth.Group("/siwe")
	siwe.Post("/nonce", authHandler.GetNonce)
	siwe.Post("/verify", authHandler.Verify)
	
	// Magic Link Authentication (stub)
	auth.Post("/magic-link", authHandler.SendMagicLink)
	
	// Get current user (protected)
	auth.Get("/me", middleware.JWTAuthWithUser(cfg.JWTSecret, userRepo), authHandler.GetMe)

	// Protected routes
	protected := v1.Use(middleware.JWTAuthWithUser(cfg.JWTSecret, userRepo))

	// Portfolio routes
	portfolio := protected.Group("/portfolio")
	portfolio.Get("/:address/balances", portfolioHandler.GetBalances)
	portfolio.Get("/:address/history", portfolioHandler.GetHistory)

	// Transaction routes
	transactions := protected.Group("/transactions")
	transactions.Get("/:address", transactionHandler.GetTransactions)
	transactions.Get("/:address/approvals", transactionHandler.GetApprovals)
	transactions.Delete("/:address/approvals/:token", transactionHandler.RevokeApproval)

	// Yield routes
	yield := protected.Group("/yield")
	
	// Pool endpoints
	yield.Get("/pools", yieldHandler.GetYieldPools)
	yield.Get("/pools/top", yieldHandler.GetTopYieldPools)
	yield.Get("/pools/protocol/:slug", yieldHandler.GetYieldPoolsByProtocol)
	yield.Get("/pools/chain/:chainId", yieldHandler.GetYieldPoolsByChain)
	
	// Position endpoints
	yield.Get("/positions/:address", yieldHandler.GetYieldPositions)
	yield.Post("/positions/:address/:positionId/claim", yieldHandler.ClaimRewards)
	
	// Protocol endpoints
	yield.Get("/protocols", yieldHandler.GetProtocols)
	
	// Admin endpoints for position management (internal use)
	yield.Post("/positions/:address", yieldHandler.CreatePosition)
	yield.Put("/positions/:positionId", yieldHandler.UpdatePosition)

	// Bridge routes
	bridge := protected.Group("/bridge")
	bridge.Post("/routes", bridgeHandler.GetBridgeRoutes)
	bridge.Post("/execute", bridgeHandler.ExecuteBridge)
	
	// Swap routes
	swap := protected.Group("/swap")
	swap.Post("/quote", swapHandler.GetSwapQuote)
	swap.Post("/execute", swapHandler.ExecuteSwap)


	// Alert routes (protected)
	alerts := protected.Group("/alerts")
	alerts.Get("/", alertHandler.GetAlerts)
	alerts.Post("/", alertHandler.CreateAlert)
	alerts.Get("/history", alertHandler.GetAlertHistory)
	alerts.Get("/:alertId", alertHandler.GetAlert)
	alerts.Patch("/:alertId", alertHandler.UpdateAlert)
	alerts.Patch("/:alertId/pause", alertHandler.PauseAlert)
	alerts.Patch("/:alertId/activate", alertHandler.ActivateAlert)
	alerts.Delete("/:alertId", alertHandler.DeleteAlert)

	// Watchlist routes (protected)
	watchlist := protected.Group("/watchlist")
	watchlist.Get("/", watchlistHandler.GetWatchlist)
	watchlist.Post("/", watchlistHandler.CreateWatchlistItem)
	watchlist.Delete("/:id", watchlistHandler.DeleteWatchlistItem)

	// Analytics routes (protected)
	analytics := protected.Group("/analytics")
	analytics.Get("/pnl/:address", analyticsHandler.GetPnL)
	analytics.Get("/export", analyticsHandler.ExportPnL)
	analytics.Get("/download", analyticsHandler.DownloadFile)
	analytics.Get("/summary/:address", analyticsHandler.GetPnLSummary)

	// Admin routes (protected + admin only)
	admin := protected.Group("/admin", middleware.AdminAuth())
	
	// User management
	admin.Get("/users", adminHandler.GetUsers)
	
	// Error logs (if available)
	admin.Get("/errors", adminHandler.GetErrors)
	
	// Feature flags
	admin.Get("/feature-flags", adminHandler.GetFeatureFlags)
	admin.Post("/feature-flags", adminHandler.CreateFeatureFlag)
	
	// System banners
	admin.Get("/banners", adminHandler.GetSystemBanners)
	admin.Post("/banners", adminHandler.CreateSystemBanner)
	admin.Put("/banners/:id", adminHandler.UpdateSystemBanner)
	admin.Delete("/banners/:id", adminHandler.DeleteSystemBanner)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return errors.NotFound("Route")
	})
}