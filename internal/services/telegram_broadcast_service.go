package services

import (
	"time"

	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/pkg/logger"
)

type TelegramBroadcastService struct {
	db             *Database
	telegramService *TelegramService
}

func NewTelegramBroadcastService(db *Database, telegramService *TelegramService) *TelegramBroadcastService {
	return &TelegramBroadcastService{
		db:              db,
		telegramService: telegramService,
	}
}

// GetTelegramBroadcasts gets all Telegram broadcasts for a user
func (s *TelegramBroadcastService) GetTelegramBroadcasts(userID uuid.UUID, status string, page int, limit int) ([]models.TelegramBroadcast, int, error) {
	var broadcasts []models.TelegramBroadcast
	var total int64

	query := s.db.DB.Where("user_id = ?", userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Model(&models.TelegramBroadcast{}).Count(&total).Error; err != nil {
		logger.Error("Failed to count Telegram broadcasts", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&broadcasts).Error; err != nil {
		logger.Error("Failed to get Telegram broadcasts", err)
		return nil, 0, err
	}

	return broadcasts, int(total), nil
}

// CreateTelegramBroadcast creates a new Telegram broadcast
func (s *TelegramBroadcastService) CreateTelegramBroadcast(name string, message string, recipients []int64, userID uuid.UUID) (*models.TelegramBroadcast, error) {
	broadcast := &models.TelegramBroadcast{
		ID:         uuid.New(),
		Name:       name,
		Message:    message,
		Recipients: recipients,
		Status:     "pending",
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.DB.Create(broadcast).Error; err != nil {
		logger.Error("Failed to create Telegram broadcast", err)
		return nil, err
	}

	return broadcast, nil
}

// GetTelegramBroadcast gets a Telegram broadcast by ID
func (s *TelegramBroadcastService) GetTelegramBroadcast(broadcastID uuid.UUID) (*models.TelegramBroadcast, error) {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast", err)
		return nil, err
	}

	return &broadcast, nil
}

// UpdateTelegramBroadcast updates a Telegram broadcast
func (s *TelegramBroadcastService) UpdateTelegramBroadcast(broadcastID uuid.UUID, name string, message string, recipients []int64) (*models.TelegramBroadcast, error) {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for update", err)
		return nil, err
	}

	// Update fields
	if name != "" {
		broadcast.Name = name
	}
	if message != "" {
		broadcast.Message = message
	}
	if recipients != nil {
		broadcast.Recipients = recipients
	}
	broadcast.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update Telegram broadcast", err)
		return nil, err
	}

	return &broadcast, nil
}

// DeleteTelegramBroadcast deletes a Telegram broadcast
func (s *TelegramBroadcastService) DeleteTelegramBroadcast(broadcastID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", broadcastID).Delete(&models.TelegramBroadcast{}).Error; err != nil {
		logger.Error("Failed to delete Telegram broadcast", err)
		return err
	}

	return nil
}

// SendTelegramBroadcast sends a Telegram broadcast
func (s *TelegramBroadcastService) SendTelegramBroadcast(broadcastID uuid.UUID) error {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for sending", err)
		return err
	}

	// Update status to sending
	broadcast.Status = "sending"
	broadcast.SentAt = time.Now()
	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update Telegram broadcast status", err)
		return err
	}

	// Send messages to recipients
	successCount := 0
	for _, recipient := range broadcast.Recipients {
		if err := s.telegramService.SendMessage(recipient, broadcast.Message); err != nil {
			logger.Error("Failed to send Telegram broadcast message", err, map[string]interface{}{
				"chat_id": recipient,
				"broadcast_id": broadcastID,
			})
		} else {
			successCount++
		}
	}

	// Update broadcast with results
	broadcast.Status = "completed"
	broadcast.SuccessCount = successCount
	broadcast.FailureCount = len(broadcast.Recipients) - successCount
	broadcast.CompletedAt = time.Now()
	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update Telegram broadcast completion", err)
		return err
	}

	return nil
}

// GetTelegramBroadcastStats gets Telegram broadcast statistics
func (s *TelegramBroadcastService) GetTelegramBroadcastStats(broadcastID uuid.UUID) (map[string]interface{}, error) {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for stats", err)
		return nil, err
	}

	stats := map[string]interface{}{
		"broadcast_id":     broadcast.ID,
		"name":             broadcast.Name,
		"status":           broadcast.Status,
		"total_recipients": len(broadcast.Recipients),
		"success_count":    broadcast.SuccessCount,
		"failure_count":    broadcast.FailureCount,
		"created_at":       broadcast.CreatedAt,
		"sent_at":          broadcast.SentAt,
		"completed_at":     broadcast.CompletedAt,
	}

	return stats, nil
}

// GetTelegramBroadcastAnalytics gets analytics for Telegram broadcasts
func (s *TelegramBroadcastService) GetTelegramBroadcastAnalytics(userID uuid.UUID, startDate time.Time, endDate time.Time) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Total broadcasts
	var totalBroadcasts int64
	s.db.DB.Model(&models.TelegramBroadcast{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Count(&totalBroadcasts)
	analytics["total_broadcasts"] = totalBroadcasts

	// Broadcasts by status
	var statusBreakdown []map[string]interface{}
	s.db.DB.Model(&models.TelegramBroadcast{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusBreakdown)
	analytics["status_breakdown"] = statusBreakdown

	// Average success rate
	var avgSuccessRate float64
	s.db.DB.Model(&models.TelegramBroadcast{}).
		Where("user_id = ? AND status = ? AND created_at BETWEEN ? AND ?", userID, "completed", startDate, endDate).
		Select("AVG(CASE WHEN total_recipients > 0 THEN (success_count * 100.0 / total_recipients) ELSE 0 END)").
		Scan(&avgSuccessRate)
	analytics["avg_success_rate"] = avgSuccessRate

	// Daily breakdown
	var dailyBreakdown []map[string]interface{}
	s.db.DB.Model(&models.TelegramBroadcast{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Select("DATE(created_at) as date, COUNT(*) as broadcasts, SUM(success_count) as total_success").
		Group("DATE(created_at)").
		Order("date").
		Scan(&dailyBreakdown)
	analytics["daily_breakdown"] = dailyBreakdown

	// Top performing broadcasts
	var topBroadcasts []models.TelegramBroadcast
	s.db.DB.Where("user_id = ? AND status = ? AND created_at BETWEEN ? AND ?", userID, "completed", startDate, endDate).
		Order("success_count desc").
		Limit(10).
		Find(&topBroadcasts)
	analytics["top_broadcasts"] = topBroadcasts

	return analytics, nil
}

// ScheduleTelegramBroadcast schedules a Telegram broadcast for future delivery
func (s *TelegramBroadcastService) ScheduleTelegramBroadcast(broadcastID uuid.UUID, scheduleAt time.Time) error {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for scheduling", err)
		return err
	}

	// In a real implementation, you would use a job scheduler like cron or a message queue
	// For now, we'll just log it and return
	logger.Info("Telegram broadcast scheduled", map[string]interface{}{
		"broadcast_id": broadcastID,
		"schedule_at":  scheduleAt,
	})

	return nil
}

// CancelTelegramBroadcast cancels a scheduled Telegram broadcast
func (s *TelegramBroadcastService) CancelTelegramBroadcast(broadcastID uuid.UUID) error {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for cancellation", err)
		return err
	}

	if broadcast.Status != "pending" {
		return fmt.Errorf("can only cancel pending broadcasts")
	}

	broadcast.Status = "cancelled"
	broadcast.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to cancel Telegram broadcast", err)
		return err
	}

	return nil
}

// DuplicateTelegramBroadcast duplicates an existing broadcast
func (s *TelegramBroadcastService) DuplicateTelegramBroadcast(broadcastID uuid.UUID, userID uuid.UUID) (*models.TelegramBroadcast, error) {
	var original models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&original).Error; err != nil {
		logger.Error("Failed to get original Telegram broadcast", err)
		return nil, err
	}

	// Create new broadcast with same content
	newBroadcast := &models.TelegramBroadcast{
		ID:         uuid.New(),
		Name:       original.Name + " (Copy)",
		Message:    original.Message,
		Recipients: original.Recipients,
		Status:     "pending",
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.DB.Create(newBroadcast).Error; err != nil {
		logger.Error("Failed to duplicate Telegram broadcast", err)
		return nil, err
	}

	return newBroadcast, nil
}

// GetTelegramBroadcastTemplates gets broadcast templates
func (s *TelegramBroadcastService) GetTelegramBroadcastTemplates(userID uuid.UUID) ([]models.TelegramBroadcast, error) {
	var templates []models.TelegramBroadcast
	
	// Get broadcasts that can be used as templates (completed broadcasts)
	if err := s.db.DB.Where("user_id = ? AND status = ?", userID, "completed").
		Order("created_at desc").
		Limit(20).
		Find(&templates).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast templates", err)
		return nil, err
	}

	return templates, nil
}

// ValidateRecipients validates Telegram recipient IDs
func (s *TelegramBroadcastService) ValidateRecipients(recipients []int64) ([]int64, []string) {
	validRecipients := make([]int64, 0)
	errors := make([]string, 0)

	for _, recipient := range recipients {
		if recipient <= 0 {
			errors = append(errors, fmt.Sprintf("Invalid recipient ID: %d", recipient))
		} else {
			validRecipients = append(validRecipients, recipient)
		}
	}

	return validRecipients, errors
}

// PreviewTelegramBroadcast shows a preview of the broadcast message
func (s *TelegramBroadcastService) PreviewTelegramBroadcast(message string, recipients []int64) (map[string]interface{}, error) {
	preview := map[string]interface{}{
		"message":           message,
		"recipient_count":   len(recipients),
		"estimated_delivery": time.Now().Add(time.Duration(len(recipients)) * time.Second),
		"message_length":    len(message),
		"has_media":         false, // In a real implementation, detect media URLs
		"has_links":         false, // In a real implementation, detect links
	}

	return preview, nil
}

// ExportTelegramBroadcastRecipients exports recipients list
func (s *TelegramBroadcastService) ExportTelegramBroadcastRecipients(broadcastID uuid.UUID) ([]int64, error) {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for export", err)
		return nil, err
	}

	return broadcast.Recipients, nil
}

// ImportTelegramBroadcastRecipients imports recipients from file
func (s *TelegramBroadcastService) ImportTelegramBroadcastRecipients(broadcastID uuid.UUID, recipients []int64) error {
	var broadcast models.TelegramBroadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get Telegram broadcast for import", err)
		return err
	}

	// Validate recipients
	validRecipients, errors := s.ValidateRecipients(recipients)
	if len(errors) > 0 {
		logger.Warn("Some recipients are invalid", map[string]interface{}{
			"errors": errors,
		})
	}

	broadcast.Recipients = validRecipients
	broadcast.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update Telegram broadcast recipients", err)
		return err
	}

	return nil
}