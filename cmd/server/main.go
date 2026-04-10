package main

import (
	"fmt"
	"log"

	"avi_bd/config"
	"avi_bd/internal/database"
	"avi_bd/internal/game"
	"avi_bd/internal/handlers"
	"avi_bd/internal/repositories"
	"avi_bd/internal/routes"
	"avi_bd/internal/services"
	"avi_bd/internal/websocket"

	_ "avi_bd/docs"

	"github.com/gin-gonic/gin"
)

// @title Aviator Betting Game API
// @version 1.0
// @description Real-time Aviator betting game backend with WebSocket support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@aviator.game

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDatabase(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create indexes
	if err := database.CreateIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Optional: seed database with test data
	// if err := database.SeedDatabase(db); err != nil {
	// 	log.Fatalf("Failed to seed database: %v", err)
	// }

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	betRepo := repositories.NewBetRepository(db)
	roundRepo := repositories.NewGameRoundRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, walletRepo, cfg)
	walletService := services.NewWalletService(walletRepo, transactionRepo)
	betService := services.NewBetService(betRepo, roundRepo, walletRepo, transactionRepo)
	profileService := services.NewProfileService(userRepo, walletRepo)
	gameService := services.NewGameService(roundRepo, betRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	// Initialize services
	adminService := services.NewAdminService(db)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize game engine
	gameEngine := game.NewEngine(db, cfg, hub)
	gameEngine.Start()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	walletHandler := handlers.NewWalletHandler(walletService)
	betHandler := handlers.NewBetHandler(betService)
	profileHandler := handlers.NewProfileHandler(profileService)
	historyHandler := handlers.NewHistoryHandler(gameService, transactionService)
	healthHandler := handlers.NewHealthHandler(db)
	adminHandler := handlers.NewAdminHandler(adminService, gameEngine)

	// Setup Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(
		router,
		cfg,
		authHandler,
		walletHandler,
		betHandler,
		profileHandler,
		historyHandler,
		healthHandler,
		adminHandler,
		hub,
	)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
