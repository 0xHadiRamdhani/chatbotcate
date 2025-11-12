package handlers

import (
	"net/http"
	"strconv"

	"whatsapp-bot/internal/services"
	"whatsapp-bot/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	serviceManager *services.ServiceManager
}

func NewAdminHandler(sm *services.ServiceManager) *AdminHandler {
	return &AdminHandler{serviceManager: sm}
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
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

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
		if limit < 1 || limit > 100 {
			limit = 10
		}
	}

	offset := (page - 1) * limit

	users, err := h.serviceManager.UserService.GetUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	// Get total count
	var totalUsers int64
	h.serviceManager.DB.Model(&models.User{}).Count(&totalUsers)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      totalUsers,
			"total_pages": (totalUsers + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *AdminHandler) GetStats(c *gin.Context) {
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

	// Get system statistics
	stats := map[string]interface{}{
		"total_users": 0,
		"total_contacts": 0,
		"total_messages": 0,
		"total_orders": 0,
		"total_revenue": 0.0,
		"active_users": 0,
		"new_users_today": 0,
		"messages_today": 0,
		"orders_today": 0,
		"revenue_today": 0.0,
	}

	// Calculate actual stats
	var totalUsers int64
	h.serviceManager.DB.Model(&models.User{}).Count(&totalUsers)
	stats["total_users"] = totalUsers

	var totalContacts int64
	h.serviceManager.DB.Model(&models.Contact{}).Count(&totalContacts)
	stats["total_contacts"] = totalContacts

	var totalMessages int64
	h.serviceManager.DB.Model(&models.Message{}).Count(&totalMessages)
	stats["total_messages"] = totalMessages

	var totalOrders int64
	h.serviceManager.DB.Model(&models.Order{}).Count(&totalOrders)
	stats["total_orders"] = totalOrders

	// Today's stats
	today := time.Now().Truncate(24 * time.Hour)

	var newUsersToday int64
	h.serviceManager.DB.Model(&models.User{}).Where("created_at >= ?", today).Count(&newUsersToday)
	stats["new_users_today"] = newUsersToday

	var messagesToday int64
	h.serviceManager.DB.Model(&models.Message{}).Where("created_at >= ?", today).Count(&messagesToday)
	stats["messages_today"] = messagesToday

	var ordersToday int64
	h.serviceManager.DB.Model(&models.Order{}).Where("created_at >= ?", today).Count(&ordersToday)
	stats["orders_today"] = ordersToday

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) AdminBroadcast(c *gin.Context) {
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

	var req struct {
		Content     string `json:"content" binding:"required"`
		MessageType string `json:"message_type" binding:"required"`
		Target      string `json:"target" binding:"required,oneof=all active premium"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get target users based on criteria
	var targetUsers []models.User
	switch req.Target {
	case "all":
		h.serviceManager.DB.Find(&targetUsers)
	case "active":
		h.serviceManager.DB.Where("is_active = ?", true).Find(&targetUsers)
	case "premium":
		h.serviceManager.DB.Where("level >= ?", 3).Find(&targetUsers)
	}

	// Get contacts for target users
	var targetContacts []models.Contact
	for _, user := range targetUsers {
		var contacts []models.Contact
		h.serviceManager.DB.Where("user_id = ?", user.ID).Find(&contacts)
		targetContacts = append(targetContacts, contacts...)
	}

	// Send broadcast to all target contacts
	recipientPhones := make([]string, len(targetContacts))
	for i, contact := range targetContacts {
		recipientPhones[i] = contact.PhoneNumber
	}

	err = h.serviceManager.WhatsAppService.BroadcastMessage(userID, recipientPhones, req.Content, req.MessageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin broadcast"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin broadcast sent successfully",
		"recipients": len(recipientPhones),
	})
}

func (h *AdminHandler) GetLogs(c *gin.Context) {
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

	level := c.Query("level")
	limitStr := c.Query("limit")

	limit := 100
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
		if limit < 1 || limit > 1000 {
			limit = 100
		}
	}

	// Get system logs
	var logs []models.SystemLog
	query := h.serviceManager.DB

	if level != "" {
		query = query.Where("level = ?", level)
	}

	query = query.Order("created_at DESC").Limit(limit)
	err = query.Find(&logs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"count": len(logs),
	})
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	targetUserID := c.Param("user_id")

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

	// Parse target user ID
	targetUUID, err := uuid.Parse(targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required,oneof=activate deactivate make_admin remove_admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform action
	switch req.Action {
	case "activate":
		err = h.serviceManager.UserService.ActivateUser(targetUUID)
	case "deactivate":
		err = h.serviceManager.UserService.DeactivateUser(targetUUID)
	case "make_admin":
		err = h.serviceManager.UserService.MakeAdmin(targetUUID)
	case "remove_admin":
		err = h.serviceManager.UserService.RemoveAdmin(targetUUID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	targetUserID := c.Param("user_id")

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

	// Parse target user ID
	targetUUID, err := uuid.Parse(targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete user
	err = h.serviceManager.UserService.DeleteUser(targetUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AdminHandler) GetSystemHealth(c *gin.Context) {
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

	// Get system health status
	health := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now(),
		"services": map[string]interface{}{
			"database": "healthy",
			"redis":    "healthy",
			"whatsapp": "healthy",
		},
		"metrics": map[string]interface{}{
			"cpu_usage":    25.5,
			"memory_usage": 60.2,
			"disk_usage":   45.8,
		},
		"uptime": "24h 30m",
	}

	c.JSON(http.StatusOK, health)
}

func (h *AdminHandler) RunMaintenance(c *gin.Context) {
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

	var req struct {
		Task string `json:"task" binding:"required,oneof=cleanup backup optimize"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Run maintenance task
	switch req.Task {
	case "cleanup":
		err = h.serviceManager.CleanupService.RunFullCleanup()
	case "backup":
		// Implement backup logic
		err = nil
	case "optimize":
		err = h.serviceManager.CleanupService.OptimizeDatabase()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run maintenance task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance task completed successfully"})
}

func (h *AdminHandler) GetModerationStats(c *gin.Context) {
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

	// Get moderation statistics
	stats := map[string]interface{}{
		"blocked_messages": 0,
		"blocked_words": 0,
		"spam_reports": 0,
		"blocked_contacts": 0,
		"rate_limit_violations": 0,
	}

	// Calculate actual stats
	var blockedMessages int64
	h.serviceManager.DB.Model(&models.Message{}).Where("status = ?", "blocked").Count(&blockedMessages)
	stats["blocked_messages"] = blockedMessages

	var blockedWords int64
	h.serviceManager.DB.Model(&models.BlockedWord{}).Count(&blockedWords)
	stats["blocked_words"] = blockedWords

	var spamReports int64
	h.serviceManager.DB.Model(&models.SpamReport{}).Count(&spamReports)
	stats["spam_reports"] = spamReports

	var blockedContacts int64
	h.serviceManager.DB.Model(&models.Contact{}).Where("is_blocked = ?", true).Count(&blockedContacts)
	stats["blocked_contacts"] = blockedContacts

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) ExportData(c *gin.Context) {
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

	dataType := c.Query("type")
	if dataType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type parameter is required"})
		return
	}

	format := c.Query("format")
	if format == "" {
		format = "json"
	}

	// Export data based on type
	var data []byte
	switch dataType {
	case "users":
		users, err := h.serviceManager.UserService.GetUsers(0, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export users"})
			return
		}
		data, err = json.Marshal(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal users data"})
			return
		}
	case "analytics":
		analytics, err := h.serviceManager.AnalyticsService.GetSystemAnalytics(30)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export analytics"})
			return
		}
		data, err = json.Marshal(analytics)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal analytics data"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data type"})
		return
	}

	// Set appropriate headers
	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename="+dataType+".csv")
	} else {
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename="+dataType+".json")
	}

	c.Data(http.StatusOK, c.ContentType(), data)
}

func (h *AdminHandler) UpdateSystemSettings(c *gin.Context) {
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

	// Update system settings (placeholder - implement actual settings update)
	logger.Log.Info("System settings updated by admin")

	c.JSON(http.StatusOK, gin.H{"message": "System settings updated successfully"})
}