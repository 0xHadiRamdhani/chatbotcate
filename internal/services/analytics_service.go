package services

import (
	"encoding/json"
	"fmt"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AnalyticsService struct {
	sm *ServiceManager
}

func (s *AnalyticsService) LogEvent(userID uuid.UUID, eventType string, value int, metadata map[string]interface{}) error {
	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	analytics := &models.Analytics{
		UserID:         userID,
		MetricType:     eventType,
		MetricValue:    value,
		AdditionalData: string(metadataJSON),
	}

	if err := s.sm.DB.Create(analytics).Error; err != nil {
		return err
	}

	return nil
}

func (s *AnalyticsService) GetUserAnalytics(userID uuid.UUID, days int) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Get date range
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Message analytics
	var messageStats struct {
		TotalSent     int64
		TotalReceived int64
		ByType        []struct {
			MessageType string
			Count       int64
		}
	}

	// Total sent messages
	s.sm.DB.Model(&models.Message{}).
		Where("user_id = ? AND direction = ? AND created_at > ?", userID, "outgoing", cutoffDate).
		Count(&messageStats.TotalSent)

	// Total received messages
	s.sm.DB.Model(&models.Message{}).
		Where("user_id = ? AND direction = ? AND created_at > ?", userID, "incoming", cutoffDate).
		Count(&messageStats.TotalReceived)

	// Messages by type
	s.sm.DB.Model(&models.Message{}).
		Where("user_id = ? AND created_at > ?", userID, cutoffDate).
		Select("message_type, count(*) as count").
		Group("message_type").
		Scan(&messageStats.ByType)

	analytics["messages"] = messageStats

	// Game analytics
	var gameStats struct {
		TotalGamesPlayed int64
		TotalScore       int64
		ByGameType       []struct {
			GameType string
			Count    int64
			Score    int64
		}
	}

	// Total games played
	s.sm.DB.Model(&models.GameScore{}).
		Where("user_id = ?", userID).
		Count(&gameStats.TotalGamesPlayed)

	// Total score
	s.sm.DB.Model(&models.GameScore{}).
		Where("user_id = ?", userID).
		Select("SUM(score) as total").
		Scan(&gameStats.TotalScore)

	// Games by type
	s.sm.DB.Model(&models.GameScore{}).
		Where("user_id = ?", userID).
		Select("game_type, COUNT(*) as count, SUM(score) as score").
		Group("game_type").
		Scan(&gameStats.ByGameType)

	analytics["games"] = gameStats

	// Business analytics (if enabled)
	var businessStats struct {
		TotalOrders     int64
		TotalRevenue    float64
		OrdersByStatus  []struct {
			Status string
			Count  int64
		}
	}

	// Total orders
	s.sm.DB.Model(&models.Order{}).
		Where("user_id = ? AND created_at > ?", userID, cutoffDate).
		Count(&businessStats.TotalOrders)

	// Total revenue
	s.sm.DB.Model(&models.Order{}).
		Where("user_id = ? AND status IN ('confirmed', 'shipped', 'delivered') AND created_at > ?", userID, cutoffDate).
		Select("SUM(total_amount) as total").
		Scan(&businessStats.TotalRevenue)

	// Orders by status
	s.sm.DB.Model(&models.Order{}).
		Where("user_id = ? AND created_at > ?", userID, cutoffDate).
		Select("status, count(*) as count").
		Group("status").
		Scan(&businessStats.OrdersByStatus)

	analytics["business"] = businessStats

	// Feature usage analytics
	var featureStats []struct {
		MetricType string
		Total      int64
	}

	s.sm.DB.Model(&models.Analytics{}).
		Where("user_id = ? AND created_at > ?", userID, cutoffDate).
		Select("metric_type, SUM(metric_value) as total").
		Group("metric_type").
		Scan(&featureStats)

	analytics["features"] = featureStats

	// Daily activity
	var dailyActivity []struct {
		Date  time.Time
		Count int64
	}

	s.sm.DB.Model(&models.Analytics{}).
		Where("user_id = ? AND created_at > ?", userID, cutoffDate).
		Select("DATE(created_at) as date, SUM(metric_value) as count").
		Group("DATE(created_at)").
		Order("date").
		Scan(&dailyActivity)

	analytics["daily_activity"] = dailyActivity

	return analytics, nil
}

func (s *AnalyticsService) GetSystemAnalytics(days int) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Get date range
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// System-wide stats
	var systemStats struct {
		TotalUsers      int64
		TotalMessages   int64
		TotalContacts   int64
		ActiveUsers     int64
		TotalOrders     int64
		TotalRevenue    float64
	}

	// Total users
	s.sm.DB.Model(&models.User{}).Count(&systemStats.TotalUsers)

	// Active users (logged in last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	s.sm.DB.Model(&models.User{}).
		Where("last_login_at > ?", thirtyDaysAgo).
		Count(&systemStats.ActiveUsers)

	// Total messages
	s.sm.DB.Model(&models.Message{}).
		Where("created_at > ?", cutoffDate).
		Count(&systemStats.TotalMessages)

	// Total contacts
	s.sm.DB.Model(&models.Contact{}).
		Where("created_at > ?", cutoffDate).
		Count(&systemStats.TotalContacts)

	// Total orders
	s.sm.DB.Model(&models.Order{}).
		Where("created_at > ?", cutoffDate).
		Count(&systemStats.TotalOrders)

	// Total revenue
	s.sm.DB.Model(&models.Order{}).
		Where("status IN ('confirmed', 'shipped', 'delivered') AND created_at > ?", cutoffDate).
		Select("SUM(total_amount) as total").
		Scan(&systemStats.TotalRevenue)

	analytics["system"] = systemStats

	// Feature popularity
	var featurePopularity []struct {
		Feature string
		Usage   int64
	}

	s.sm.DB.Model(&models.Analytics{}).
		Where("created_at > ?", cutoffDate).
		Select("metric_type as feature, SUM(metric_value) as usage").
		Group("metric_type").
		Order("usage DESC").
		Limit(10).
		Scan(&featurePopularity)

	analytics["feature_popularity"] = featurePopularity

	// User growth
	var userGrowth []struct {
		Date  time.Time
		Count int64
	}

	s.sm.DB.Model(&models.User{}).
		Where("created_at > ?", cutoffDate).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Order("date").
		Scan(&userGrowth)

	analytics["user_growth"] = userGrowth

	// Message trends
	var messageTrends []struct {
		Date  time.Time
		Sent  int64
		Received int64
	}

	s.sm.DB.Model(&models.Message{}).
		Where("created_at > ?", cutoffDate).
		Select("DATE(created_at) as date, SUM(CASE WHEN direction = 'outgoing' THEN 1 ELSE 0 END) as sent, SUM(CASE WHEN direction = 'incoming' THEN 1 ELSE 0 END) as received").
		Group("DATE(created_at)").
		Order("date").
		Scan(&messageTrends)

	analytics["message_trends"] = messageTrends

	return analytics, nil
}

func (s *AnalyticsService) GetFeatureAnalytics(feature string, days int) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Get date range
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Feature usage over time
	var usageOverTime []struct {
		Date  time.Time
		Count int64
	}

	s.sm.DB.Model(&models.Analytics{}).
		Where("metric_type = ? AND created_at > ?", feature, cutoffDate).
		Select("DATE(created_at) as date, SUM(metric_value) as count").
		Group("DATE(created_at)").
		Order("date").
		Scan(&usageOverTime)

	analytics["usage_over_time"] = usageOverTime

	// Top users for this feature
	var topUsers []struct {
		UserID uuid.UUID
		Usage  int64
	}

	s.sm.DB.Model(&models.Analytics{}).
		Where("metric_type = ? AND created_at > ?", feature, cutoffDate).
		Select("user_id, SUM(metric_value) as usage").
		Group("user_id").
		Order("usage DESC").
		Limit(10).
		Scan(&topUsers)

	analytics["top_users"] = topUsers

	// Average usage per user
	var avgUsage float64
	s.sm.DB.Model(&models.Analytics{}).
		Where("metric_type = ? AND created_at > ?", feature, cutoffDate).
		Select("AVG(metric_value) as avg").
		Scan(&avgUsage)

	analytics["average_usage"] = avgUsage

	return analytics, nil
}

func (s *AnalyticsService) GenerateReport(userID uuid.UUID, reportType string, startDate, endDate time.Time) ([]byte, error) {
	// Generate comprehensive report
	report := make(map[string]interface{})

	// Report metadata
	report["metadata"] = map[string]interface{}{
		"type":       reportType,
		"start_date": startDate,
		"end_date":   endDate,
		"generated":  time.Now(),
	}

	// Get analytics data based on report type
	switch reportType {
	case "user_summary":
		analytics, err := s.GetUserAnalytics(userID, 30)
		if err != nil {
			return nil, err
		}
		report["data"] = analytics

	case "business_summary":
		businessStats, err := s.sm.BusinessService.GetBusinessStats(userID)
		if err != nil {
			return nil, err
		}
		report["data"] = businessStats

	case "game_summary":
		gameStats, err := s.sm.GameService.GetLeaderboard("", 0)
		if err != nil {
			return nil, err
		}
		report["data"] = gameStats

	default:
		return nil, fmt.Errorf("invalid report type")
	}

	// Convert to JSON
	return json.Marshal(report)
}

func (s *AnalyticsService) ExportAnalytics(userID uuid.UUID, format string, startDate, endDate time.Time) ([]byte, error) {
	// Get analytics data
	analytics, err := s.GetUserAnalytics(userID, 30)
	if err != nil {
		return nil, err
	}

	switch format {
	case "json":
		return json.Marshal(analytics)

	case "csv":
		// Convert to CSV format
		return s.convertToCSV(analytics)

	default:
		return nil, fmt.Errorf("unsupported format")
	}
}

func (s *AnalyticsService) convertToCSV(data map[string]interface{}) ([]byte, error) {
	// Simple CSV conversion
	var csvData []byte

	// Add headers
	csvData = append(csvData, []byte("Metric,Value\n"))

	// Add data rows
	for key, value := range data {
		row := fmt.Sprintf("%s,%v\n", key, value)
		csvData = append(csvData, []byte(row))
	}

	return csvData, nil
}

func (s *AnalyticsService) CleanUpOldAnalytics(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	return s.sm.DB.Where("created_at < ?", cutoffDate).
		Delete(&models.Analytics{}).Error
}

func (s *AnalyticsService) GetRealTimeStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Messages in last hour
	var messagesLastHour int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	s.sm.DB.Model(&models.Message{}).
		Where("created_at > ?", oneHourAgo).
		Count(&messagesLastHour)
	stats["messages_last_hour"] = messagesLastHour

	// Active users (last 5 minutes)
	var activeUsers int64
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	s.sm.DB.Model(&models.User{}).
		Where("last_login_at > ?", fiveMinutesAgo).
		Count(&activeUsers)
	stats["active_users_now"] = activeUsers

	// New users today
	var newUsersToday int64
	today := time.Now().Truncate(24 * time.Hour)
	s.sm.DB.Model(&models.User{}).
		Where("created_at > ?", today).
		Count(&newUsersToday)
	stats["new_users_today"] = newUsersToday

	// System load (messages per minute)
	var messagesPerMinute float64
	s.sm.DB.Model(&models.Message{}).
		Where("created_at > ?", time.Now().Add(-1*time.Minute)).
		Select("COUNT(*) / 1.0 as rate").
		Scan(&messagesPerMinute)
	stats["messages_per_minute"] = messagesPerMinute

	return stats, nil
}

func (s *AnalyticsService) TrackUserJourney(userID uuid.UUID) ([]UserJourneyStep, error) {
	var journey []UserJourneyStep

	// Get user events in chronological order
	var events []models.Analytics
	err := s.sm.DB.Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	// Process events into journey steps
	for _, event := range events {
		step := UserJourneyStep{
			Timestamp: event.CreatedAt,
			Action:    event.MetricType,
			Value:     event.MetricValue,
		}

		// Parse metadata
		if event.AdditionalData != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(event.AdditionalData), &metadata); err == nil {
				step.Metadata = metadata
			}
		}

		journey = append(journey, step)
	}

	return journey, nil
}

type UserJourneyStep struct {
	Timestamp time.Time
	Action    string
	Value     int
	Metadata  map[string]interface{}
}

func (s *AnalyticsService) PredictUserBehavior(userID uuid.UUID) (map[string]interface{}, error) {
	// Simple behavior prediction based on historical data
	prediction := make(map[string]interface{})

	// Get recent activity
	recentEvents, err := s.TrackUserJourney(userID)
	if err != nil {
		return nil, err
	}

	// Analyze patterns
	gameUsage := 0
	businessUsage := 0
	utilityUsage := 0

	for _, event := range recentEvents {
		switch event.Action {
		case "game_played", "quiz_completed", "khodam_checked":
			gameUsage++
		case "order_created", "catalog_viewed", "payment_confirmed":
			businessUsage++
		case "weather_checked", "currency_converted", "qr_code_generated":
			utilityUsage++
		}
	}

	// Predict next likely action
	if gameUsage > businessUsage && gameUsage > utilityUsage {
		prediction["likely_next_action"] = "game"
		prediction["confidence"] = float64(gameUsage) / float64(len(recentEvents))
	} else if businessUsage > utilityUsage {
		prediction["likely_next_action"] = "business"
		prediction["confidence"] = float64(businessUsage) / float64(len(recentEvents))
	} else {
		prediction["likely_next_action"] = "utility"
		prediction["confidence"] = float64(utilityUsage) / float64(len(recentEvents))
	}

	return prediction, nil
}

func (s *AnalyticsService) CreateCustomDashboard(userID uuid.UUID, name string, widgets []DashboardWidget) error {
	// Store custom dashboard configuration
	dashboard := map[string]interface{}{
		"user_id":   userID,
		"name":      name,
		"widgets":   widgets,
		"created_at": time.Now(),
	}

	// Store in database or cache
	logger.Log.WithFields(logger.Fields{
		"user_id": userID,
		"name":    name,
		"widgets": len(widgets),
	}).Info("Custom dashboard created")

	return nil
}

type DashboardWidget struct {
	Type   string                 `json:"type"`
	Title  string                 `json:"title"`
	Config map[string]interface{} `json:"config"`
}