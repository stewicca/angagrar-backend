package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stewicca/angagrar-backend/config"
	"github.com/stewicca/angagrar-backend/internal/database"
	"github.com/stewicca/angagrar-backend/internal/handlers"
	"github.com/stewicca/angagrar-backend/internal/middleware"
	"github.com/stewicca/angagrar-backend/internal/repositories"
	"github.com/stewicca/angagrar-backend/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	db := database.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	budgetRepo := repositories.NewBudgetRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	openAIService := services.NewOpenAIService(cfg)
	conversationService := services.NewConversationService(
		conversationRepo,
		messageRepo,
		budgetRepo,
		openAIService,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	budgetHandler := handlers.NewBudgetHandler(budgetRepo)
	conversationHandler := handlers.NewConversationHandler(conversationService)

	r := gin.Default()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Angagrar Backend API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/guest", authHandler.CreateGuest)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			users.GET("/profile", userHandler.GetProfile)
		}

		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetTransactions)
		}

		conversations := api.Group("/conversations")
		conversations.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			conversations.POST("/start", conversationHandler.StartConversation)
			conversations.POST("/:sessionId/messages", conversationHandler.SendMessage)
			conversations.GET("/:sessionId/history", conversationHandler.GetConversationHistory)
			conversations.POST("/:sessionId/reset", conversationHandler.ResetConversation)
		}

		budgets := api.Group("/budgets")
		budgets.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			budgets.GET("", budgetHandler.GetUserBudgets)
			budgets.PATCH("/:id", budgetHandler.UpdateBudget)
		}
	}

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
