package handlers

import (
	"time"

	"kilocode.dev/whatsapp-bot/internal/middleware"
	"kilocode.dev/whatsapp-bot/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, serviceManager *services.ServiceManager) {
	// Public routes
	public := router.Group("/api/v1")
	{
		// Auth routes
		authHandler := NewAuthHandler(serviceManager.UserService)
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)
		public.POST("/auth/forgot-password", authHandler.ForgotPassword)
		public.POST("/auth/reset-password", authHandler.ResetPassword)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.RateLimitMiddleware())
	{
		// User routes
		userHandler := NewUserHandler(serviceManager.UserService)
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.DELETE("/users/profile", userHandler.DeleteProfile)
		protected.POST("/users/change-password", userHandler.ChangePassword)
		protected.GET("/users/settings", userHandler.GetSettings)
		protected.PUT("/users/settings", userHandler.UpdateSettings)

		// WhatsApp routes
		whatsappHandler := NewWhatsAppHandler(serviceManager.WhatsAppService)
		protected.POST("/whatsapp/send-message", whatsappHandler.SendMessage)
		protected.POST("/whatsapp/send-media", whatsappHandler.SendMedia)
		protected.POST("/whatsapp/send-template", whatsappHandler.SendTemplate)
		protected.GET("/whatsapp/contacts", whatsappHandler.GetContacts)
		protected.GET("/whatsapp/chats", whatsappHandler.GetChats)
		protected.GET("/whatsapp/chat/:chat_id", whatsappHandler.GetChatMessages)
		protected.POST("/whatsapp/mark-read", whatsappHandler.MarkAsRead)
		protected.GET("/whatsapp/status", whatsappHandler.GetStatus)

		// Auto-reply routes
		autoReplyHandler := NewAutoReplyHandler(serviceManager.AutoReplyService)
		protected.GET("/auto-replies", autoReplyHandler.GetAutoReplies)
		protected.POST("/auto-replies", autoReplyHandler.CreateAutoReply)
		protected.GET("/auto-replies/:reply_id", autoReplyHandler.GetAutoReply)
		protected.PUT("/auto-replies/:reply_id", autoReplyHandler.UpdateAutoReply)
		protected.DELETE("/auto-replies/:reply_id", autoReplyHandler.DeleteAutoReply)
		protected.POST("/auto-replies/:reply_id/toggle", autoReplyHandler.ToggleAutoReply)

		// Broadcast routes
		broadcastHandler := NewBroadcastHandler(serviceManager.BroadcastService)
		protected.GET("/broadcasts", broadcastHandler.GetBroadcasts)
		protected.POST("/broadcasts", broadcastHandler.CreateBroadcast)
		protected.GET("/broadcasts/:broadcast_id", broadcastHandler.GetBroadcast)
		protected.PUT("/broadcasts/:broadcast_id", broadcastHandler.UpdateBroadcast)
		protected.DELETE("/broadcasts/:broadcast_id", broadcastHandler.DeleteBroadcast)
		protected.POST("/broadcasts/:broadcast_id/send", broadcastHandler.SendBroadcast)
		protected.GET("/broadcasts/:broadcast_id/stats", broadcastHandler.GetBroadcastStats)

		// Game routes
		gameHandler := NewGameHandler(serviceManager.GameService)
		protected.GET("/games", gameHandler.GetGames)
		protected.POST("/games/start", gameHandler.StartGame)
		protected.POST("/games/:game_id/play", gameHandler.PlayGame)
		protected.GET("/games/:game_id/leaderboard", gameHandler.GetLeaderboard)
		protected.GET("/games/:game_id/stats", gameHandler.GetGameStats)
		protected.POST("/games/trivia/start", gameHandler.StartTrivia)
		protected.POST("/games/trivia/answer", gameHandler.AnswerTrivia)

		// Utility routes
		utilityHandler := NewUtilityHandler(serviceManager.UtilityService)
		protected.POST("/utils/qr-code", utilityHandler.CreateQRCode)
		protected.POST("/utils/short-link", utilityHandler.CreateShortLink)
		protected.GET("/utils/short-link/:alias", utilityHandler.GetShortLink)
		protected.POST("/utils/currency-convert", utilityHandler.ConvertCurrency)
		protected.GET("/utils/weather", utilityHandler.GetWeather)
		protected.POST("/utils/translate", utilityHandler.TranslateText)
		protected.POST("/utils/location-info", utilityHandler.GetLocationInfo)
		protected.POST("/utils/polls", utilityHandler.CreatePoll)
		protected.POST("/utils/polls/:poll_id/vote", utilityHandler.VotePoll)
		protected.GET("/utils/polls/:poll_id/results", utilityHandler.GetPollResults)
		protected.POST("/utils/reminders", utilityHandler.CreateReminder)
		protected.GET("/utils/reminders", utilityHandler.GetReminders)
		protected.DELETE("/utils/reminders/:reminder_id", utilityHandler.DeleteReminder)
		protected.POST("/utils/notes", utilityHandler.CreateNote)
		protected.GET("/utils/notes", utilityHandler.GetNotes)
		protected.PUT("/utils/notes/:note_id", utilityHandler.UpdateNote)
		protected.DELETE("/utils/notes/:note_id", utilityHandler.DeleteNote)
		protected.GET("/utils/notes/search", utilityHandler.SearchNotes)
		protected.POST("/utils/timers", utilityHandler.CreateTimer)
		protected.GET("/utils/timers/:timer_id", utilityHandler.GetTimer)
		protected.DELETE("/utils/timers/:timer_id", utilityHandler.StopTimer)
		protected.POST("/utils/file-upload", utilityHandler.CreateFileUpload)
		protected.GET("/utils/file-upload/:upload_id", utilityHandler.GetFileUpload)
		protected.DELETE("/utils/file-upload/:upload_id", utilityHandler.DeleteFileUpload)

		// Business routes
		businessHandler := NewBusinessHandler(serviceManager.BusinessService)
		protected.POST("/business/products", businessHandler.CreateProduct)
		protected.GET("/business/products", businessHandler.GetProducts)
		protected.GET("/business/products/:product_id", businessHandler.GetProduct)
		protected.PUT("/business/products/:product_id", businessHandler.UpdateProduct)
		protected.DELETE("/business/products/:product_id", businessHandler.DeleteProduct)
		protected.POST("/business/orders", businessHandler.CreateOrder)
		protected.GET("/business/orders", businessHandler.GetOrders)
		protected.GET("/business/orders/:order_id", businessHandler.GetOrder)
		protected.PUT("/business/orders/:order_id/status", businessHandler.UpdateOrderStatus)
		protected.POST("/business/invoices", businessHandler.CreateInvoice)
		protected.GET("/business/invoices", businessHandler.GetInvoices)
		protected.GET("/business/invoices/:invoice_id", businessHandler.GetInvoice)
		protected.PUT("/business/invoices/:invoice_id/status", businessHandler.UpdateInvoiceStatus)
		protected.POST("/business/customers", businessHandler.CreateCustomer)
		protected.GET("/business/customers", businessHandler.GetCustomers)
		protected.GET("/business/customers/:customer_id", businessHandler.GetCustomer)
		protected.PUT("/business/customers/:customer_id", businessHandler.UpdateCustomer)
		protected.DELETE("/business/customers/:customer_id", businessHandler.DeleteCustomer)
		protected.POST("/business/payments", businessHandler.CreatePayment)
		protected.GET("/business/payments", businessHandler.GetPayments)
		protected.GET("/business/payments/:payment_id", businessHandler.GetPayment)
		protected.GET("/business/stats", businessHandler.GetBusinessStats)
		protected.GET("/business/sales-report", businessHandler.GetSalesReport)

		// Moderation routes
		moderationHandler := NewModerationHandler(serviceManager.ModerationService)
		protected.GET("/moderation/blocked-users", moderationHandler.GetBlockedUsers)
		protected.POST("/moderation/block-user", moderationHandler.BlockUser)
		protected.DELETE("/moderation/unblock-user/:user_id", moderationHandler.UnblockUser)
		protected.GET("/moderation/reported-messages", moderationHandler.GetReportedMessages)
		protected.POST("/moderation/report-message", moderationHandler.ReportMessage)
		protected.PUT("/moderation/reported-messages/:report_id", moderationHandler.UpdateReportStatus)
		protected.GET("/moderation/spam-detections", moderationHandler.GetSpamDetections)
		protected.POST("/moderation/mark-spam", moderationHandler.MarkAsSpam)
		protected.GET("/moderation/content-filters", moderationHandler.GetContentFilters)
		protected.POST("/moderation/content-filters", moderationHandler.CreateContentFilter)
		protected.PUT("/moderation/content-filters/:filter_id", moderationHandler.UpdateContentFilter)
		protected.DELETE("/moderation/content-filters/:filter_id", moderationHandler.DeleteContentFilter)

		// Analytics routes
		analyticsHandler := NewAnalyticsHandler(serviceManager.AnalyticsService)
		protected.GET("/analytics/dashboard", analyticsHandler.GetDashboard)
		protected.GET("/analytics/messages", analyticsHandler.GetMessageAnalytics)
		protected.GET("/analytics/users", analyticsHandler.GetUserAnalytics)
		protected.GET("/analytics/games", analyticsHandler.GetGameAnalytics)
		protected.GET("/analytics/business", analyticsHandler.GetBusinessAnalytics)
		protected.GET("/analytics/export", analyticsHandler.ExportAnalytics)
	}

	// Admin routes (requires admin role)
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		adminHandler := NewAdminHandler(serviceManager)
		admin.GET("/dashboard", adminHandler.GetDashboard)
		admin.GET("/users", adminHandler.GetUsers)
		admin.GET("/users/:user_id", adminHandler.GetUser)
		admin.PUT("/users/:user_id", adminHandler.UpdateUser)
		admin.DELETE("/users/:user_id", adminHandler.DeleteUser)
		admin.POST("/users/:user_id/ban", adminHandler.BanUser)
		admin.POST("/users/:user_id/unban", adminHandler.UnbanUser)
		admin.GET("/system-stats", adminHandler.GetSystemStats)
		admin.GET("/logs", adminHandler.GetLogs)
		admin.DELETE("/logs", adminHandler.ClearLogs)
		admin.GET("/settings", adminHandler.GetSettings)
		admin.PUT("/settings", adminHandler.UpdateSettings)
		admin.POST("/backup", adminHandler.CreateBackup)
		admin.GET("/backups", adminHandler.GetBackups)
		admin.POST("/backups/:backup_id/restore", adminHandler.RestoreBackup)
		admin.DELETE("/backups/:backup_id", adminHandler.DeleteBackup)
		admin.GET("/system-health", adminHandler.GetSystemHealth)
		admin.POST("/system-maintenance", adminHandler.SetMaintenanceMode)
		admin.GET("/broadcast-messages", adminHandler.GetBroadcastMessages)
		admin.POST("/broadcast-messages", adminHandler.CreateBroadcastMessage)
		admin.DELETE("/broadcast-messages/:message_id", adminHandler.DeleteBroadcastMessage)
		admin.GET("/spam-reports", adminHandler.GetSpamReports)
		admin.PUT("/spam-reports/:report_id", adminHandler.UpdateSpamReport)
		admin.GET("/content-reports", adminHandler.GetContentReports)
		admin.PUT("/content-reports/:report_id", adminHandler.UpdateContentReport)
	}

	// Health check route (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// Webhook routes for WhatsApp
	webhook := router.Group("/webhooks")
	{
		whatsappHandler := NewWhatsAppHandler(serviceManager.WhatsAppService)
		webhook.POST("/whatsapp", whatsappHandler.HandleWebhook)
		webhook.GET("/whatsapp", whatsappHandler.VerifyWebhook)
	}
}