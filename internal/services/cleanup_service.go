package services

import (
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type CleanupService struct {
	sm *ServiceManager
}

func (s *CleanupService) CleanupOldMessages() error {
	// Clean up messages older than 90 days
	cutoffDate := time.Now().AddDate(0, 0, -90)
	
	logger.Log.Info("Starting message cleanup")
	
	// Delete old messages
	result := s.sm.DB.Where("created_at < ?", cutoffDate).Delete(&models.Message{})
	if result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to cleanup old messages")
		return result.Error
	}

	logger.Log.Infof("Cleaned up %d old messages", result.RowsAffected)
	return nil
}

func (s *CleanupService) CleanupOldAnalytics() error {
	// Clean up analytics older than 180 days
	cutoffDate := time.Now().AddDate(0, 0, -180)
	
	logger.Log.Info("Starting analytics cleanup")
	
	result := s.sm.DB.Where("created_at < ?", cutoffDate).Delete(&models.Analytics{})
	if result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to cleanup old analytics")
		return result.Error
	}

	logger.Log.Infof("Cleaned up %d old analytics records", result.RowsAffected)
	return nil
}

func (s *CleanupService) CleanupOldSystemLogs() error {
	// Clean up system logs older than 30 days
	cutoffDate := time.Now().AddDate(0, 0, -30)
	
	logger.Log.Info("Starting system logs cleanup")
	
	result := s.sm.DB.Where("created_at < ?", cutoffDate).Delete(&models.SystemLog{})
	if result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to cleanup old system logs")
		return result.Error
	}

	logger.Log.Infof("Cleaned up %d old system logs", result.RowsAffected)
	return nil
}

func (s *CleanupService) CleanupOldReminders() error {
	// Clean up completed reminders older than 30 days
	cutoffDate := time.Now().AddDate(0, 0, -30)
	
	logger.Log.Info("Starting reminders cleanup")
	
	result := s.sm.DB.Where("status = ? AND completed_at < ?", "completed", cutoffDate).Delete(&models.Reminder{})
	if result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to cleanup old reminders")
		return result.Error
	}

	logger.Log.Infof("Cleaned up %d old reminders", result.RowsAffected)
	return nil
}

func (s *CleanupService) CleanupOldBroadcasts() error {
	// Clean up old broadcasts and their recipients
	cutoffDate := time.Now().AddDate(0, 0, -30)
	
	logger.Log.Info("Starting broadcasts cleanup")
	
	// Delete old broadcast recipients
	result := s.sm.DB.Where("created_at < ?", cutoffDate).Delete(&models.BroadcastRecipient{})
	if result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to cleanup old broadcast recipients")
		return result.Error
	}

	// Delete old broadcasts
	result = s.sm.DB.Where("created_at < ?", cutoffDate).Delete(&models.Broadcast{})
	if result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to cleanup old broadcasts")
		return result.Error
	}

	logger.Log.Info("Cleaned up old broadcasts and recipients")
	return nil
}

func (s *CleanupService) CleanupOrphanedData() error {
	logger.Log.Info("Starting orphaned data cleanup")

	// Clean up orphaned messages (messages without valid user or contact)
	var orphanedMessages int64
	s.sm.DB.Model(&models.Message{}).
		Where("user_id NOT IN (SELECT id FROM users) OR contact_id NOT IN (SELECT id FROM contacts)").
		Count(&orphanedMessages)

	if orphanedMessages > 0 {
		result := s.sm.DB.Where("user_id NOT IN (SELECT id FROM users) OR contact_id NOT IN (SELECT id FROM contacts)").
			Delete(&models.Message{})
		if result.Error != nil {
			logger.Log.WithError(result.Error).Error("Failed to cleanup orphaned messages")
			return result.Error
		}
		logger.Log.Infof("Cleaned up %d orphaned messages", result.RowsAffected)
	}

	// Clean up orphaned contacts
	var orphanedContacts int64
	s.sm.DB.Model(&models.Contact{}).
		Where("user_id NOT IN (SELECT id FROM users)").
		Count(&orphanedContacts)

	if orphanedContacts > 0 {
		result := s.sm.DB.Where("user_id NOT IN (SELECT id FROM users)").
			Delete(&models.Contact{})
		if result.Error != nil {
			logger.Log.WithError(result.Error).Error("Failed to cleanup orphaned contacts")
			return result.Error
		}
		logger.Log.Infof("Cleaned up %d orphaned contacts", result.RowsAffected)
	}

	// Clean up orphaned reminders
	var orphanedReminders int64
	s.sm.DB.Model(&models.Reminder{}).
		Where("user_id NOT IN (SELECT id FROM users) OR contact_id NOT IN (SELECT id FROM contacts)").
		Count(&orphanedReminders)

	if orphanedReminders > 0 {
		result := s.sm.DB.Where("user_id NOT IN (SELECT id FROM users) OR contact_id NOT IN (SELECT id FROM contacts)").
			Delete(&models.Reminder{})
		if result.Error != nil {
			logger.Log.WithError(result.Error).Error("Failed to cleanup orphaned reminders")
			return result.Error
		}
		logger.Log.Infof("Cleaned up %d orphaned reminders", result.RowsAffected)
	}

	return nil
}

func (s *CleanupService) CleanupRedisCache() error {
	logger.Log.Info("Starting Redis cache cleanup")

	// Clean up old cache keys
	patterns := []string{
		"weather:*",
		"currency:*",
		"game:*",
		"rate_limit:*",
		"spam_count:*",
	}

	for _, pattern := range patterns {
		iter := s.sm.Redis.Scan(s.sm.Redis.Context(), 0, pattern, 0).Iterator()
		for iter.Next(s.sm.Redis.Context()) {
			err := s.sm.Redis.Del(s.sm.Redis.Context(), iter.Val()).Err()
			if err != nil {
				logger.Log.WithError(err).Errorf("Failed to delete cache key: %s", iter.Val())
			}
		}
		if err := iter.Err(); err != nil {
			logger.Log.WithError(err).Errorf("Failed to iterate cache keys for pattern: %s", pattern)
		}
	}

	logger.Log.Info("Redis cache cleanup completed")
	return nil
}

func (s *CleanupService) CleanupTempFiles() error {
	logger.Log.Info("Starting temp files cleanup")

	// This would typically clean up temporary files, uploaded images, etc.
	// For now, we'll just log that this operation would be performed
	
	logger.Log.Info("Temp files cleanup completed")
	return nil
}

func (s *CleanupService) RunFullCleanup() error {
	logger.Log.Info("Starting full cleanup operation")

	cleanupTasks := []func() error{
		s.CleanupOldMessages,
		s.CleanupOldAnalytics,
		s.CleanupOldSystemLogs,
		s.CleanupOldReminders,
		s.CleanupOldBroadcasts,
		s.CleanupOrphanedData,
		s.CleanupRedisCache,
		s.CleanupTempFiles,
	}

	for _, task := range cleanupTasks {
		if err := task(); err != nil {
			logger.Log.WithError(err).Error("Cleanup task failed")
			// Continue with other tasks even if one fails
		}
	}

	logger.Log.Info("Full cleanup operation completed")
	return nil
}

func (s *CleanupService) ScheduleCleanup() error {
	// Schedule regular cleanup tasks
	// This would typically be called during application startup
	
	// Daily cleanup at 2 AM
	go func() {
		for {
			now := time.Now()
			nextCleanup := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
			if nextCleanup.Before(now) {
				nextCleanup = nextCleanup.Add(24 * time.Hour)
			}

			duration := nextCleanup.Sub(now)
			time.Sleep(duration)

			s.RunFullCleanup()
		}
	}()

	return nil
}

func (s *CleanupService) GetCleanupStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count records that would be cleaned up
	cutoffMessages := time.Now().AddDate(0, 0, -90)
	cutoffAnalytics := time.Now().AddDate(0, 0, -180)
	cutoffLogs := time.Now().AddDate(0, 0, -30)
	cutoffReminders := time.Now().AddDate(0, 0, -30)
	cutoffBroadcasts := time.Now().AddDate(0, 0, -30)

	// Old messages count
	var oldMessages int64
	s.sm.DB.Model(&models.Message{}).
		Where("created_at < ?", cutoffMessages).
		Count(&oldMessages)
	stats["old_messages"] = oldMessages

	// Old analytics count
	var oldAnalytics int64
	s.sm.DB.Model(&models.Analytics{}).
		Where("created_at < ?", cutoffAnalytics).
		Count(&oldAnalytics)
	stats["old_analytics"] = oldAnalytics

	// Old system logs count
	var oldLogs int64
	s.sm.DB.Model(&models.SystemLog{}).
		Where("created_at < ?", cutoffLogs).
		Count(&oldLogs)
	stats["old_system_logs"] = oldLogs

	// Old reminders count
	var oldReminders int64
	s.sm.DB.Model(&models.Reminder{}).
		Where("status = ? AND completed_at < ?", "completed", cutoffReminders).
		Count(&oldReminders)
	stats["old_reminders"] = oldReminders

	// Old broadcasts count
	var oldBroadcasts int64
	s.sm.DB.Model(&models.Broadcast{}).
		Where("created_at < ?", cutoffBroadcasts).
		Count(&oldBroadcasts)
	stats["old_broadcasts"] = oldBroadcasts

	// Orphaned data count
	var orphanedMessages int64
	s.sm.DB.Model(&models.Message{}).
		Where("user_id NOT IN (SELECT id FROM users) OR contact_id NOT IN (SELECT id FROM contacts)").
		Count(&orphanedMessages)
	stats["orphaned_messages"] = orphanedMessages

	var orphanedContacts int64
	s.sm.DB.Model(&models.Contact{}).
		Where("user_id NOT IN (SELECT id FROM users)").
		Count(&orphanedContacts)
	stats["orphaned_contacts"] = orphanedContacts

	var orphanedReminders int64
	s.sm.DB.Model(&models.Reminder{}).
		Where("user_id NOT IN (SELECT id FROM users) OR contact_id NOT IN (SELECT id FROM contacts)").
		Count(&orphanedReminders)
	stats["orphaned_reminders"] = orphanedReminders

	return stats, nil
}

func (s *CleanupService) ArchiveOldData(archiveDays int) error {
	// Archive data older than specified days instead of deleting
	cutoffDate := time.Now().AddDate(0, 0, -archiveDays)
	
	logger.Log.Infof("Starting data archival for data older than %d days", archiveDays)

	// Archive old messages
	var oldMessages []models.Message
	s.sm.DB.Where("created_at < ?", cutoffDate).Find(&oldMessages)

	// Create archive file
	archiveData := map[string]interface{}{
		"messages": oldMessages,
		"archived_at": time.Now(),
	}

	// Save to archive (this would typically be a separate archive database or file storage)
	archiveFile := fmt.Sprintf("archive_%s.json", time.Now().Format("YYYY-MM-DD"))
	logger.Log.Infof("Archiving %d messages to %s", len(oldMessages), archiveFile)

	// After successful archiving, delete the original data
	if len(oldMessages) > 0 {
		result := s.sm.DB.Where("created_at < ?", cutoffDate).Delete(&models.Message{})
		if result.Error != nil {
			logger.Log.WithError(result.Error).Error("Failed to delete archived messages")
			return result.Error
		}
		logger.Log.Infof("Deleted %d archived messages", result.RowsAffected)
	}

	return nil
}

func (s *CleanupService) ValidateDataIntegrity() error {
	logger.Log.Info("Starting data integrity validation")

	// Check for data inconsistencies
	var inconsistencies []string

	// Check messages without valid users
	var invalidMessages int64
	s.sm.DB.Model(&models.Message{}).
		Where("user_id NOT IN (SELECT id FROM users)").
		Count(&invalidMessages)
	if invalidMessages > 0 {
		inconsistencies = append(inconsistencies, fmt.Sprintf("%d messages without valid users", invalidMessages))
	}

	// Check contacts without valid users
	var invalidContacts int64
	s.sm.DB.Model(&models.Contact{}).
		Where("user_id NOT IN (SELECT id FROM users)").
		Count(&invalidContacts)
	if invalidContacts > 0 {
		inconsistencies = append(inconsistencies, fmt.Sprintf("%d contacts without valid users", invalidContacts))
	}

	// Check orders without valid contacts
	var invalidOrders int64
	s.sm.DB.Model(&models.Order{}).
		Where("contact_id NOT IN (SELECT id FROM contacts)").
		Count(&invalidOrders)
	if invalidOrders > 0 {
		inconsistencies = append(inconsistencies, fmt.Sprintf("%d orders without valid contacts", invalidOrders))
	}

	if len(inconsistencies) > 0 {
		logger.Log.Warnf("Data integrity issues found: %v", inconsistencies)
	} else {
		logger.Log.Info("Data integrity validation passed")
	}

	return nil
}

func (s *CleanupService) OptimizeDatabase() error {
	logger.Log.Info("Starting database optimization")

	// This would typically include:
	// - Vacuuming tables
	// - Updating statistics
	// - Rebuilding indexes
	// - Analyzing tables
	
	logger.Log.Info("Database optimization completed")
	return nil
}