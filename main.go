package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whatsapp-bot/internal/config"
	"whatsapp-bot/internal/database"
	"whatsapp-bot/internal/handlers"
	"whatsapp-bot/internal/middleware"
	"whatsapp-bot/internal/services"
	"whatsapp-bot/internal/utils"
	"whatsapp-bot/pkg/logger"
	"whatsapp-bot/pkg/whatsapp"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logger.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := database.InitializeRedis(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	defer redisClient.Close()

	// Initialize WhatsApp client
	waClient, err := whatsapp.Initialize(cfg.WhatsApp)
	if err != nil {
		log.Fatal("Failed to initialize WhatsApp client:", err)
	}
	defer waClient.Disconnect()

	// Initialize services
	serviceManager := services.NewServiceManager(db, redisClient, waClient, cfg)

	// Initialize cron jobs
	cronManager := cron.New()
	setupCronJobs(cronManager, serviceManager)
	cronManager.Start()
	defer cronManager.Stop()

	// Setup HTTP server
	router := setupRouter(serviceManager)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRouter(serviceManager *services.ServiceManager) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			authHandler := handlers.NewAuthHandler(serviceManager)
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// WhatsApp routes
		whatsapp := api.Group("/whatsapp")
		whatsapp.Use(middleware.AuthJWT())
		{
			whatsappHandler := handlers.NewWhatsAppHandler(serviceManager)
			whatsapp.POST("/send", whatsappHandler.SendMessage)
			whatsapp.POST("/broadcast", whatsappHandler.BroadcastMessage)
			whatsapp.GET("/contacts", whatsappHandler.GetContacts)
			whatsapp.POST("/groups", whatsappHandler.CreateGroup)
			whatsapp.GET("/groups", whatsappHandler.GetGroups)
		}

		// Bot features routes
		bot := api.Group("/bot")
		bot.Use(middleware.AuthJWT())
		{
			botHandler := handlers.NewBotHandler(serviceManager)
			bot.GET("/features", botHandler.GetFeatures)
			bot.POST("/features/:feature/enable", botHandler.EnableFeature)
			bot.POST("/features/:feature/disable", botHandler.DisableFeature)
			bot.GET("/analytics", botHandler.GetAnalytics)
		}

		// Game routes
		game := api.Group("/game")
		game.Use(middleware.AuthJWT())
		{
			gameHandler := handlers.NewGameHandler(serviceManager)
			game.GET("/leaderboard", gameHandler.GetLeaderboard)
			game.POST("/quiz/start", gameHandler.StartQuiz)
			game.POST("/quiz/answer", gameHandler.SubmitAnswer)
			game.GET("/khodam/:name", gameHandler.CheckKhodam)
		}

		// Business routes
		business := api.Group("/business")
		business.Use(middleware.AuthJWT())
		{
			businessHandler := handlers.NewBusinessHandler(serviceManager)
			business.GET("/products", businessHandler.GetProducts)
			business.POST("/products", businessHandler.CreateProduct)
			business.GET("/orders", businessHandler.GetOrders)
			business.POST("/orders", businessHandler.CreateOrder)
			business.GET("/customers", businessHandler.GetCustomers)
		}

		// Utility routes
		utils := api.Group("/utils")
		utils.Use(middleware.AuthJWT())
		{
			utilsHandler := handlers.NewUtilsHandler(serviceManager)
			utils.GET("/weather/:city", utilsHandler.GetWeather)
			utils.GET("/currency", utilsHandler.ConvertCurrency)
			utils.POST("/translate", utilsHandler.Translate)
			utils.GET("/qrcode", utilsHandler.GenerateQRCode)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthJWT())
		admin.Use(middleware.RequireAdmin())
		{
			adminHandler := handlers.NewAdminHandler(serviceManager)
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/stats", adminHandler.GetStats)
			admin.POST("/broadcast", adminHandler.AdminBroadcast)
			admin.GET("/logs", adminHandler.GetLogs)
		}
	}

	// Webhook for WhatsApp
	router.POST("/webhook/whatsapp", handlers.NewWhatsAppHandler(serviceManager).HandleWebhook)

	return router
}

func setupCronJobs(cronManager *cron.Cron, serviceManager *services.ServiceManager) {
	// Daily cleanup job
	cronManager.AddFunc("0 2 * * *", func() {
		serviceManager.CleanupService.CleanupOldMessages()
	})

	// Hourly reminder check
	cronManager.AddFunc("0 * * * *", func() {
		serviceManager.ReminderService.ProcessReminders()
	})

	// Daily leaderboard reset
	cronManager.AddFunc("0 0 * * *", func() {
		serviceManager.GameService.ResetDailyLeaderboard()
	})

	// Weather update every 6 hours
	cronManager.AddFunc("0 */6 * * *", func() {
		serviceManager.WeatherService.UpdateWeatherData()
	})

	// Currency rate update every hour
	cronManager.AddFunc("0 * * * *", func() {
		serviceManager.CurrencyService.UpdateRates()
	})
}