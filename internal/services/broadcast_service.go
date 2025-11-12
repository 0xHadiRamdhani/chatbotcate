package services

import (
	"time"

	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/pkg/logger"
)

type BroadcastService struct {
	db *Database
}

func NewBroadcastService(db *Database) *BroadcastService {
	return &BroadcastService{
		db: db,
	}
}

// GetBroadcasts gets all broadcasts for a user
func (s *BroadcastService) GetBroadcasts(userID uuid.UUID, status string, page int, limit int) ([]models.Broadcast, int, error) {
	var broadcasts []models.Broadcast
	var total int64

	query := s.db.DB.Where("user_id = ?", userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Model(&models.Broadcast{}).Count(&total).Error; err != nil {
		logger.Error("Failed to count broadcasts", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&broadcasts).Error; err != nil {
		logger.Error("Failed to get broadcasts", err)
		return nil, 0, err
	}

	return broadcasts, int(total), nil
}

// CreateBroadcast creates a new broadcast
func (s *BroadcastService) CreateBroadcast(name string, message string, recipients []string, scheduleAt string, userID uuid.UUID) (*models.Broadcast, error) {
	broadcast := &models.Broadcast{
		ID:         uuid.New(),
		Name:       name,
		Message:    message,
		Recipients: recipients,
		Status:     "pending",
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Parse schedule time if provided
	if scheduleAt != "" {
		scheduleTime, err := time.Parse("YYYY-MM-DD HH:mm:ss", scheduleAt)
		if err != nil {
			logger.Error("Invalid schedule time format", err)
			return nil, err
		}
		broadcast.ScheduleAt = &scheduleTime
	}

	if err := s.db.DB.Create(broadcast).Error; err != nil {
		logger.Error("Failed to create broadcast", err)
		return nil, err
	}

	return broadcast, nil
}

// GetBroadcast gets a broadcast by ID
func (s *BroadcastService) GetBroadcast(broadcastID uuid.UUID) (*models.Broadcast, error) {
	var broadcast models.Broadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get broadcast", err)
		return nil, err
	}

	return &broadcast, nil
}

// UpdateBroadcast updates a broadcast
func (s *BroadcastService) UpdateBroadcast(broadcastID uuid.UUID, name string, message string, recipients []string, scheduleAt string) (*models.Broadcast, error) {
	var broadcast models.Broadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get broadcast for update", err)
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
	if scheduleAt != "" {
		scheduleTime, err := time.Parse("YYYY-MM-DD HH:mm:ss", scheduleAt)
		if err != nil {
			logger.Error("Invalid schedule time format", err)
			return nil, err
		}
		broadcast.ScheduleAt = &scheduleTime
	}
	broadcast.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update broadcast", err)
		return nil, err
	}

	return &broadcast, nil
}

// DeleteBroadcast deletes a broadcast
func (s *BroadcastService) DeleteBroadcast(broadcastID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", broadcastID).Delete(&models.Broadcast{}).Error; err != nil {
		logger.Error("Failed to delete broadcast", err)
		return err
	}

	return nil
}

// SendBroadcast sends a broadcast
func (s *BroadcastService) SendBroadcast(broadcastID uuid.UUID) error {
	var broadcast models.Broadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get broadcast for sending", err)
		return err
	}

	// Update status to sending
	broadcast.Status = "sending"
	broadcast.SentAt = time.Now()
	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update broadcast status", err)
		return err
	}

	// Send messages to recipients
	successCount := 0
	for _, recipient := range broadcast.Recipients {
		// Here you would integrate with your WhatsApp service
		// For now, we'll just simulate success
		successCount++
	}

	// Update broadcast with results
	broadcast.Status = "completed"
	broadcast.SuccessCount = successCount
	broadcast.FailureCount = len(broadcast.Recipients) - successCount
	broadcast.CompletedAt = time.Now()
	if err := s.db.DB.Save(&broadcast).Error; err != nil {
		logger.Error("Failed to update broadcast completion", err)
		return err
	}

	return nil
}

// GetBroadcastStats gets broadcast statistics
func (s *BroadcastService) GetBroadcastStats(broadcastID uuid.UUID) (map[string]interface{}, error) {
	var broadcast models.Broadcast
	if err := s.db.DB.Where("id = ?", broadcastID).First(&broadcast).Error; err != nil {
		logger.Error("Failed to get broadcast for stats", err)
		return nil, err
	}

	stats := map[string]interface{}{
		"broadcast_id":  broadcast.ID,
		"name":          broadcast.Name,
		"status":        broadcast.Status,
		"total_recipients": len(broadcast.Recipients),
		"success_count": broadcast.SuccessCount,
		"failure_count": broadcast.FailureCount,
		"created_at":    broadcast.CreatedAt,
		"sent_at":       broadcast.SentAt,
		"completed_at":  broadcast.CompletedAt,
	}

	return stats, nil
}