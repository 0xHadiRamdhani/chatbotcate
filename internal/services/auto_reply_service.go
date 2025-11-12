package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AutoReplyService struct {
	sm *ServiceManager
}

func (s *AutoReplyService) ProcessAutoReply(contact *models.Contact, message *models.Message) error {
	// Get active auto-replies for user
	var autoReplies []models.AutoReply
	err := s.sm.DB.Where("user_id = ? AND is_active = ?", contact.UserID, true).Find(&autoReplies).Error
	if err != nil {
		return err
	}

	for _, autoReply := range autoReplies {
		if s.shouldTriggerAutoReply(message.Content, autoReply) {
			// Send auto-reply
			if err := s.sendAutoReply(contact, autoReply); err != nil {
				logger.Log.WithError(err).Error("Failed to send auto-reply")
				continue
			}

			// Log analytics
			s.sm.AnalyticsService.LogEvent(contact.UserID, "auto_reply_sent", 1, map[string]interface{}{
				"keyword": autoReply.Keyword,
			})

			return nil // Only trigger first matching auto-reply
		}
	}

	return nil
}

func (s *AutoReplyService) shouldTriggerAutoReply(message string, autoReply models.AutoReply) bool {
	message = strings.ToLower(strings.TrimSpace(message))
	keyword := strings.ToLower(strings.TrimSpace(autoReply.Keyword))

	switch autoReply.MatchType {
	case "exact":
		return message == keyword
	case "contains":
		return strings.Contains(message, keyword)
	case "regex":
		matched, err := regexp.MatchString(keyword, message)
		return err == nil && matched
	default:
		return false
	}
}

func (s *AutoReplyService) sendAutoReply(contact *models.Contact, autoReply models.AutoReply) error {
	switch autoReply.ReplyType {
	case "text":
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, autoReply.Response, false)
		return err
	case "image":
		if autoReply.MediaURL != "" {
			_, err := s.sm.WhatsApp.SendImageMessage(contact.PhoneNumber, autoReply.MediaURL, autoReply.Response)
			return err
		}
		// Fallback to text if no image URL
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, autoReply.Response, false)
		return err
	case "template":
		// Implement template message
		return fmt.Errorf("template messages not implemented yet")
	default:
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, autoReply.Response, false)
		return err
	}
}

func (s *AutoReplyService) CreateAutoReply(userID uuid.UUID, keyword, response, matchType, replyType string, mediaURL string) (*models.AutoReply, error) {
	autoReply := &models.AutoReply{
		UserID:    userID,
		Keyword:   keyword,
		Response:  response,
		MatchType: matchType,
		ReplyType: replyType,
		MediaURL:  mediaURL,
		IsActive:  true,
	}

	if err := s.sm.DB.Create(autoReply).Error; err != nil {
		return nil, err
	}

	return autoReply, nil
}

func (s *AutoReplyService) GetAutoReplies(userID uuid.UUID) ([]models.AutoReply, error) {
	var autoReplies []models.AutoReply
	err := s.sm.DB.Where("user_id = ?", userID).Find(&autoReplies).Error
	return autoReplies, err
}

func (s *AutoReplyService) UpdateAutoReply(id uuid.UUID, updates map[string]interface{}) error {
	return s.sm.DB.Model(&models.AutoReply{}).Where("id = ?", id).Updates(updates).Error
}

func (s *AutoReplyService) DeleteAutoReply(id uuid.UUID) error {
	return s.sm.DB.Where("id = ?", id).Delete(&models.AutoReply{}).Error
}

func (s *AutoReplyService) ToggleAutoReply(id uuid.UUID) error {
	var autoReply models.AutoReply
	if err := s.sm.DB.Where("id = ?", id).First(&autoReply).Error; err != nil {
		return err
	}

	autoReply.IsActive = !autoReply.IsActive
	return s.sm.DB.Save(&autoReply).Error
}

func (s *AutoReplyService) GetAutoReplyStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total auto-replies
	var total int64
	s.sm.DB.Model(&models.AutoReply{}).Where("user_id = ?", userID).Count(&total)
	stats["total_auto_replies"] = total

	// Active auto-replies
	var active int64
	s.sm.DB.Model(&models.AutoReply{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&active)
	stats["active_auto_replies"] = active

	// Auto-replies by match type
	var byMatchType []struct {
		MatchType string
		Count     int64
	}
	s.sm.DB.Model(&models.AutoReply{}).
		Where("user_id = ?", userID).
		Select("match_type, count(*) as count").
		Group("match_type").
		Scan(&byMatchType)
	stats["by_match_type"] = byMatchType

	return stats, nil
}

func (s *AutoReplyService) CreateWelcomeMessage(userID uuid.UUID, message string) (*models.AutoReply, error) {
	return s.CreateAutoReply(userID, "welcome", message, "exact", "text", "")
}

func (s *AutoReplyService) CreateAwayMessage(userID uuid.UUID, message string) (*models.AutoReply, error) {
	return s.CreateAutoReply(userID, "away", message, "exact", "text", "")
}

func (s *AutoReplyService) CreateBusinessHoursReply(userID uuid.UUID, businessHours string) (*models.AutoReply, error) {
	message := fmt.Sprintf("🕐 Jam operasional kami:\n%s\n\nKami akan segera membalas pesan Anda saat jam kerja.", businessHours)
	return s.CreateAutoReply(userID, "jam", message, "contains", "text", "")
}

func (s *AutoReplyService) CreateFAQReplies(userID uuid.UUID, faqs map[string]string) error {
	for question, answer := range faqs {
		keyword := fmt.Sprintf("faq_%s", strings.ToLower(strings.ReplaceAll(question, " ", "_")))
		if _, err := s.CreateAutoReply(userID, keyword, answer, "contains", "text", ""); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoReplyService) CreateQuickReplyButtons(userID uuid.UUID, keyword string, message string, buttons []string) error {
	// Create interactive message with buttons
	interactiveButtons := make([]whatsapp.ReplyButton, len(buttons))
	for i, button := range buttons {
		interactiveButtons[i] = whatsapp.ReplyButton{
			Type: "reply",
			Reply: whatsapp.Reply{
				ID:    fmt.Sprintf("button_%d", i),
				Title: button,
			},
		}
	}

	// For now, we'll store this as a template and handle it in the processing
	autoReply := &models.AutoReply{
		UserID:      userID,
		Keyword:     keyword,
		Response:    message,
		MatchType:   "exact",
		ReplyType:   "interactive",
		IsActive:    true,
	}

	return s.sm.DB.Create(autoReply).Error
}

func (s *AutoReplyService) ProcessInteractiveResponse(userID uuid.UUID, messageID, buttonID string) error {
	// Handle button responses
	buttonIndex := 0
	if strings.HasPrefix(buttonID, "button_") {
		indexStr := strings.TrimPrefix(buttonID, "button_")
		if index, err := fmt.Sscanf(indexStr, "%d", &buttonIndex); err == nil && index == 1 {
			// Process button response
			response := fmt.Sprintf("Anda memilih opsi %d", buttonIndex+1)
			_, err := s.sm.WhatsApp.SendTextMessage("recipient_phone", response, false)
			return err
		}
	}

	return nil
}

func (s *AutoReplyService) CreateTimeBasedAutoReply(userID uuid.UUID, keyword, response string, startTime, endTime time.Time) (*models.AutoReply, error) {
	// Check if current time is within the specified time range
	now := time.Now()
	if now.After(startTime) && now.Before(endTime) {
		return s.CreateAutoReply(userID, keyword, response, "exact", "text", "")
	}

	return nil, fmt.Errorf("outside specified time range")
}

func (s *AutoReplyService) CreateConditionalAutoReply(userID uuid.UUID, keyword string, conditions map[string]string, responses map[string]string) error {
	// Create conditional auto-reply based on multiple criteria
	// This is a placeholder for more complex conditional logic
	for condition, response := range responses {
		if _, err := s.CreateAutoReply(userID, fmt.Sprintf("%s_%s", keyword, condition), response, "contains", "text", ""); err != nil {
			return err
		}
	}

	return nil
}

func (s *AutoReplyService) GetAutoReplyByID(id uuid.UUID) (*models.AutoReply, error) {
	var autoReply models.AutoReply
	err := s.sm.DB.Where("id = ?", id).First(&autoReply).Error
	return &autoReply, err
}

func (s *AutoReplyService) SearchAutoReplies(userID uuid.UUID, query string) ([]models.AutoReply, error) {
	var autoReplies []models.AutoReply
	err := s.sm.DB.Where("user_id = ? AND (keyword LIKE ? OR response LIKE ?)", userID, "%"+query+"%", "%"+query+"%").Find(&autoReplies).Error
	return autoReplies, err
}

func (s *AutoReplyService) BulkCreateAutoReplies(userID uuid.UUID, autoReplies []models.AutoReply) error {
	tx := s.sm.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, autoReply := range autoReplies {
		autoReply.UserID = userID
		if err := tx.Create(&autoReply).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *AutoReplyService) ExportAutoReplies(userID uuid.UUID) ([]byte, error) {
	autoReplies, err := s.GetAutoReplies(userID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(autoReplies)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *AutoReplyService) ImportAutoReplies(userID uuid.UUID, data []byte) error {
	var autoReplies []models.AutoReply
	if err := json.Unmarshal(data, &autoReplies); err != nil {
		return err
	}

	return s.BulkCreateAutoReplies(userID, autoReplies)
}