package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboard gets analytics dashboard
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	userID := c.Query("user_id")
	timeRange := c.Query("time_range")

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	// Parse time range
	var startDate, endDate time.Time
	switch timeRange {
	case "today":
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
		endDate = time.Now()
	default:
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	}

	dashboard, err := h.analyticsService.GetDashboard(userUUID, startDate, endDate)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, dashboard)
}

// GetMessageAnalytics gets message analytics
func (h *AnalyticsHandler) GetMessageAnalytics(c *gin.Context) {
	userID := c.Query("user_id")
	timeRange := c.Query("time_range")
	groupBy := c.Query("group_by")

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	// Parse time range
	var startDate, endDate time.Time
	switch timeRange {
	case "today":
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
		endDate = time.Now()
	default:
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	}

	analytics, err := h.analyticsService.GetMessageAnalytics(userUUID, startDate, endDate, groupBy)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, analytics)
}

// GetUserAnalytics gets user analytics
func (h *AnalyticsHandler) GetUserAnalytics(c *gin.Context) {
	timeRange := c.Query("time_range")
	groupBy := c.Query("group_by")

	// Parse time range
	var startDate, endDate time.Time
	switch timeRange {
	case "today":
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
		endDate = time.Now()
	default:
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	}

	analytics, err := h.analyticsService.GetUserAnalytics(startDate, endDate, groupBy)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, analytics)
}

// GetGameAnalytics gets game analytics
func (h *AnalyticsHandler) GetGameAnalytics(c *gin.Context) {
	timeRange := c.Query("time_range")
	gameType := c.Query("game_type")

	// Parse time range
	var startDate, endDate time.Time
	switch timeRange {
	case "today":
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
		endDate = time.Now()
	default:
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	}

	analytics, err := h.analyticsService.GetGameAnalytics(startDate, endDate, gameType)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, analytics)
}

// GetBusinessAnalytics gets business analytics
func (h *AnalyticsHandler) GetBusinessAnalytics(c *gin.Context) {
	userID := c.Query("user_id")
	timeRange := c.Query("time_range")

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	// Parse time range
	var startDate, endDate time.Time
	switch timeRange {
	case "today":
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
		endDate = time.Now()
	default:
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	}

	analytics, err := h.analyticsService.GetBusinessAnalytics(userUUID, startDate, endDate)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, analytics)
}

// ExportAnalytics exports analytics data
func (h *AnalyticsHandler) ExportAnalytics(c *gin.Context) {
	userID := c.Query("user_id")
	timeRange := c.Query("time_range")
	format := c.Query("format")

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	// Parse time range
	var startDate, endDate time.Time
	switch timeRange {
	case "today":
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
		endDate = time.Now()
	default:
		startDate = time.Now().AddDate(0, 0, -7)
		endDate = time.Now()
	}

	// Default to CSV format
	if format == "" {
		format = "csv"
	}

	exportData, err := h.analyticsService.ExportAnalytics(userUUID, startDate, endDate, format)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set appropriate headers for file download
	filename := "analytics_export_" + time.Now().Format("YYYY-MM-DD") + "." + format
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", exportData)
}