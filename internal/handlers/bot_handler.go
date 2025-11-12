package handlers

import (
	"net/http"

	"whatsapp-bot/internal/services"
	"whatsapp-bot/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BotHandler struct {
	serviceManager *services.ServiceManager
}

func NewBotHandler(sm *services.ServiceManager) *BotHandler {
	return &BotHandler{serviceManager: sm}
}

func (h *BotHandler) GetFeatures(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Get user preferences to check enabled features
	preferences, err := h.serviceManager.UserService.GetUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user preferences"})
		return
	}

	features := map[string]bool{
		"games":      preferences.EnableGames,
		"business":   preferences.EnableBusiness,
		"utils":      preferences.EnableUtils,
		"moderation": preferences.EnableModeration,
	}

	c.JSON(http.StatusOK, gin.H{
		"features": features,
	})
}

func (h *BotHandler) EnableFeature(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	feature := c.Param("feature")

	// Validate feature
	allowedFeatures := []string{"games", "business", "utils", "moderation"}
	if !contains(allowedFeatures, feature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feature"})
		return
	}

	// Update user preferences
	updates := map[string]interface{}{
		"enable_" + feature: true,
	}

	err := h.serviceManager.UserService.UpdateUserPreferences(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable feature"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feature enabled successfully"})
}

func (h *BotHandler) DisableFeature(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	feature := c.Param("feature")

	// Validate feature
	allowedFeatures := []string{"games", "business", "utils", "moderation"}
	if !contains(allowedFeatures, feature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feature"})
		return
	}

	// Update user preferences
	updates := map[string]interface{}{
		"enable_" + feature: false,
	}

	err := h.serviceManager.UserService.UpdateUserPreferences(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable feature"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feature disabled successfully"})
}

func (h *BotHandler) GetAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Get analytics data
	analytics, err := h.serviceManager.AnalyticsService.GetUserAnalytics(userID, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (h *BotHandler) GetSystemStats(c *gin.Context) {
	// Get system-wide analytics
	analytics, err := h.serviceManager.AnalyticsService.GetSystemAnalytics(30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (h *BotHandler) GetRealTimeStats(c *gin.Context) {
	// Get real-time statistics
	stats, err := h.serviceManager.AnalyticsService.GetRealTimeStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get real-time stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *BotHandler) ExportAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	format := c.Query("format")
	if format == "" {
		format = "json"
	}

	// Validate format
	allowedFormats := []string{"json", "csv"}
	if !contains(allowedFormats, format) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Allowed: json, csv"})
		return
	}

	// Export analytics
	data, err := h.serviceManager.AnalyticsService.ExportAnalytics(userID, format, time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export analytics"})
		return
	}

	// Set appropriate headers
	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=analytics.csv")
	} else {
		c.Header("Content-Type", "application/json")
	}

	c.Data(http.StatusOK, c.ContentType(), data)
}

func (h *BotHandler) GenerateReport(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		ReportType string    `json:"report_type" binding:"required"`
		StartDate  time.Time `json:"start_date"`
		EndDate    time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default dates if not provided
	if req.StartDate.IsZero() {
		req.StartDate = time.Now().AddDate(0, 0, -30)
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now()
	}

	// Generate report
	report, err := h.serviceManager.AnalyticsService.GenerateReport(userID, req.ReportType, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	c.Data(http.StatusOK, "application/json", report)
}

func (h *BotHandler) GetUserJourney(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Get user journey
	journey, err := h.serviceManager.AnalyticsService.TrackUserJourney(userID)
	if err != nil {
		c.JSON(http.StatusInternalSpace, gin.H{"error": "Failed to get user journey"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"journey": journey,
	})
}

func (h *BotHandler) PredictUserBehavior(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Predict user behavior
	prediction, err := h.serviceManager.AnalyticsService.PredictUserBehavior(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to predict user behavior"})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

func (h *BotHandler) CreateCustomDashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		Name    string                 `json:"name" binding:"required"`
		Widgets []DashboardWidget      `json:"widgets" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create custom dashboard
	err := h.serviceManager.AnalyticsService.CreateCustomDashboard(userID, req.Name, req.Widgets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create custom dashboard"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Custom dashboard created successfully"})
}

func (h *BotHandler) GetBotStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Get bot status
	status := map[string]interface{}{
		"status":      "online",
		"uptime":      time.Since(time.Now().Add(-time.Hour)).String(), // Placeholder
		"total_users": 0, // Placeholder
		"features": map[string]bool{
			"games":      true,
			"business":   true,
			"utils":      true,
			"moderation": true,
		},
		"version": "1.0.0",
	}

	c.JSON(http.StatusOK, status)
}

func (h *BotHandler) RestartBot(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Check if user is admin
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Restart bot (placeholder - implement actual restart logic)
	logger.Log.Info("Bot restart requested by admin")

	c.JSON(http.StatusOK, gin.H{"message": "Bot restart initiated"})
}

func (h *BotHandler) GetBotLogs(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Check if user is admin
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Get recent logs (placeholder - implement actual log retrieval)
	logs := []map[string]interface{}{
		{
			"timestamp": time.Now(),
			"level":     "info",
			"message":   "Bot started successfully",
		},
		{
			"timestamp": time.Now().Add(-time.Minute),
			"level":     "warn",
			"message":   "High memory usage detected",
		},
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *BotHandler) UpdateBotSettings(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Check if user is admin
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update bot settings (placeholder - implement actual settings update)
	logger.Log.Info("Bot settings updated by admin")

	c.JSON(http.StatusOK, gin.H{"message": "Bot settings updated successfully"})
}

func (h *BotHandler) GetCleanupStats(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Check if user is admin
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Get cleanup statistics
	stats, err := h.serviceManager.CleanupService.GetCleanupStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cleanup stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *BotHandler) RunCleanup(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	// Check if user is admin
	user, err := h.serviceManager.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Run cleanup
	err = h.serviceManager.CleanupService.RunFullCleanup()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run cleanup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cleanup completed successfully"})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type DashboardWidget struct {
	Type   string                 `json:"type"`
	Title  string                 `json:"title"`
	Config map[string]interface{} `json:"config"`
}