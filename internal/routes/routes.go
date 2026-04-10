// package routes

// import (
// 	"net/http"

// 	"avi_bd/config"
// 	"avi_bd/internal/auth"
// 	"avi_bd/internal/handlers"
// 	"avi_bd/internal/websocket"

// 	"github.com/gin-gonic/gin"
// 	gorillawebsocket "github.com/gorilla/websocket"
// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // SetupRoutes sets up all routes
// func SetupRoutes(
// 	router *gin.Engine,
// 	cfg *config.Config,
// 	authHandler *handlers.AuthHandler,
// 	walletHandler *handlers.WalletHandler,
// 	betHandler *handlers.BetHandler,
// 	profileHandler *handlers.ProfileHandler,
// 	historyHandler *handlers.HistoryHandler,
// 	healthHandler *handlers.HealthHandler,
// 	hub *websocket.Hub,
// ) {
// 	// Public routes
// 	public := router.Group("/api/v1")
// 	{
// 		// Health check
// 		public.GET("/health", healthHandler.Health)

// 		// Authentication
// 		public.POST("/register", authHandler.Register)
// 		public.POST("/login", authHandler.Login)
// 	}

// 	// Protected routes
// 	protected := router.Group("/api/v1")
// 	protected.Use(auth.JWTMiddleware(cfg))
// 	{
// 		// Authentication
// 		protected.GET("/verify", authHandler.Verify)

// 		// Profile
// 		protected.GET("/profile", profileHandler.GetProfile)
// 		protected.PUT("/profile", profileHandler.UpdateProfile)
// 		protected.GET("/profile/statistics", profileHandler.GetStatistics)

// 		// Wallet
// 		protected.GET("/wallet", walletHandler.GetWallet)

// 		// Betting
// 		protected.POST("/bet", betHandler.PlaceBet)
// 		protected.POST("/cashout", betHandler.Cashout)

// 		// History
// 		protected.GET("/bets/history", betHandler.GetBetHistory)
// 		protected.GET("/game/history", historyHandler.GetGameHistory)
// 		protected.GET("/transactions/history", historyHandler.GetTransactionHistory)
// 	}

// 	// Swagger documentation
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// WebSocket
// 	router.GET("/ws", handleWebSocket(hub))
// }

// // handleWebSocket handles WebSocket connections
// func handleWebSocket(hub *websocket.Hub) gin.HandlerFunc {
// 	var upgrader = gorillawebsocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true // Allow all origins in development; restrict in production
// 		},
// 	}

// 	return func(c *gin.Context) {
// 		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 		if err != nil {
// 			return
// 		}

// 		client := websocket.NewClient(hub, conn)
// 		hub.Register(client)

// 		go client.ReadLoop()
// 		go client.WriteLoop()
// 	}
// }

package routes

import (
	"net/http"

	"avi_bd/config"
	"avi_bd/internal/auth"
	"avi_bd/internal/handlers"
	"avi_bd/internal/websocket"

	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// corsMiddleware allows all origins — safe for local development
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// SetupRoutes sets up all routes
func SetupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	walletHandler *handlers.WalletHandler,
	betHandler *handlers.BetHandler,
	profileHandler *handlers.ProfileHandler,
	historyHandler *handlers.HistoryHandler,
	healthHandler *handlers.HealthHandler,
	adminHandler *handlers.AdminHandler,
	hub *websocket.Hub,
) {
	// Allow cross-origin requests from Flutter Web (localhost:3000 / localhost:3001)
	router.Use(corsMiddleware())

	// WebSocket — at root so Flutter can reach it at ws://host:8080/ws
	router.GET("/ws", handleWebSocket(hub))

	// Public routes
	public := router.Group("/api/v1")
	{
		// Health check
		public.GET("/health", healthHandler.Health)

		// Authentication
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(auth.JWTMiddleware(cfg))
	{
		// Authentication
		protected.GET("/verify", authHandler.Verify)

		// Profile
		protected.GET("/profile", profileHandler.GetProfile)
		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.GET("/profile/statistics", profileHandler.GetStatistics)

		// Wallet
		protected.GET("/wallet", walletHandler.GetWallet)

		// Betting
		protected.POST("/bet", betHandler.PlaceBet)
		protected.POST("/cashout", betHandler.Cashout)

		// History
		protected.GET("/bets/history", betHandler.GetBetHistory)
		protected.GET("/game/history", historyHandler.GetGameHistory)
		protected.GET("/transactions/history", historyHandler.GetTransactionHistory)

		// Admin / Analytics
		admin := protected.Group("/admin")
		{
			admin.GET("/stats", adminHandler.GetStats)
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/users/:id/bets", adminHandler.GetUserBets)
			admin.GET("/crash-mode", adminHandler.GetCrashMode)
			admin.POST("/crash-mode", adminHandler.SetCrashMode)
			admin.POST("/force-crash", adminHandler.ForceCrash)
		}

		// Wallet deposit / withdraw (simulated — no real payment)
		protected.POST("/wallet/deposit", walletHandler.Deposit)
		protected.POST("/wallet/withdraw", walletHandler.Withdraw)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// handleWebSocket handles WebSocket connections
func handleWebSocket(hub *websocket.Hub) gin.HandlerFunc {
	upgrader := gorillawebsocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development; restrict in production
		},
	}

	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "WebSocket upgrade failed",
				"error":   err.Error(),
			})
			return
		}

		client := websocket.NewClient(hub, conn)
		hub.Register(client)

		go client.ReadLoop()
		go client.WriteLoop()
	}
}
