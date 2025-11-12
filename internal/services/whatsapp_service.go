package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"
	"whatsapp-bot/pkg/whatsapp"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type WhatsAppService struct {
	sm *ServiceManager
}

func (s *WhatsAppService) SendMessage(userID uuid.UUID, to, content, messageType string) (*models.Message, error) {
	// Create WhatsApp message request
	var waReq whatsapp.MessageRequest
	waReq.MessagingProduct = "whatsapp"
	waReq.RecipientType = "individual"
	waReq.To = to

	switch messageType {
	case "text":
		waReq.Type = "text"
		waReq.Text = &whatsapp.TextMessage{
			Body: content,
		}
	case "image":
		waReq.Type = "image"
		waReq.Image = &whatsapp.MediaMessage{
			Link: content,
		}
	default:
		waReq.Type = "text"
		waReq.Text = &whatsapp.TextMessage{
			Body: content,
		}
	}

	// Send message via WhatsApp API
	resp, err := s.sm.WhatsApp.SendMessage(waReq)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to send WhatsApp message")
		return nil, err
	}

	// Save message to database
	message := &models.Message{
		UserID:       userID,
		MessageID:    resp.Messages[0].ID,
		Content:      content,
		MessageType:  messageType,
		Direction:    "outgoing",
		Status:       "sent",
		Timestamp:    time.Now(),
	}

	if err := s.sm.DB.Create(message).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to save message to database")
		return nil, err
	}

	// Update contact last message time
	s.sm.ContactService.UpdateLastMessageTime(to, time.Now())

	return message, nil
}

func (s *WhatsAppService) BroadcastMessage(userID uuid.UUID, recipients []string, content, messageType string) error {
	broadcast := &models.Broadcast{
		UserID:      userID,
		Name:        fmt.Sprintf("Broadcast_%d", time.Now().Unix()),
		Content:     content,
		MessageType: messageType,
		Status:      "sending",
		SentAt:      &time.Time{},
	}

	if err := s.sm.DB.Create(broadcast).Error; err != nil {
		return err
	}

	// Process recipients in batches
	batchSize := 50
	for i := 0; i < len(recipients); i += batchSize {
		end := i + batchSize
		if end > len(recipients) {
			end = len(recipients)
		}

		batch := recipients[i:end]
		if err := s.processBroadcastBatch(broadcast, batch); err != nil {
			logger.Log.WithError(err).Error("Failed to process broadcast batch")
		}
	}

	// Update broadcast status
	broadcast.Status = "sent"
	broadcast.SentAt = &time.Now{}
	return s.sm.DB.Save(broadcast).Error
}

func (s *WhatsAppService) processBroadcastBatch(broadcast *models.Broadcast, recipients []string) error {
	for _, recipient := range recipients {
		recipientModel := &models.BroadcastRecipient{
			BroadcastID: broadcast.ID,
			Status:      "pending",
		}

		// Find or create contact
		contact, err := s.sm.ContactService.FindOrCreateContact(recipient)
		if err != nil {
			recipientModel.Status = "failed"
			recipientModel.Error = err.Error()
			s.sm.DB.Create(recipientModel)
			continue
		}

		recipientModel.ContactID = contact.ID

		// Send message
		message, err := s.SendMessage(broadcast.UserID, recipient, broadcast.Content, broadcast.MessageType)
		if err != nil {
			recipientModel.Status = "failed"
			recipientModel.Error = err.Error()
			broadcast.TotalFailed++
		} else {
			recipientModel.Status = "sent"
			recipientModel.SentAt = &time.Now()
			recipientModel.MessageID = message.ID
			broadcast.TotalSent++
		}

		s.sm.DB.Create(recipientModel)
	}

	return nil
}

func (s *WhatsAppService) HandleIncomingMessage(payload *whatsapp.WebhookPayload) error {
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			for _, message := range change.Value.Messages {
				if err := s.processIncomingMessage(&message, change.Value.Contacts); err != nil {
					logger.Log.WithError(err).Error("Failed to process incoming message")
				}
			}
		}
	}
	return nil
}

func (s *WhatsAppService) processIncomingMessage(message *whatsapp.Message, contacts []whatsapp.Contact) error {
	// Find sender contact
	var senderContact *whatsapp.Contact
	for _, contact := range contacts {
		if contact.WaID == message.From {
			senderContact = &contact
			break
		}
	}

	if senderContact == nil {
		return fmt.Errorf("sender contact not found")
	}

	// Find or create contact in database
	contact, err := s.sm.ContactService.FindOrCreateContact(message.From)
	if err != nil {
		return err
	}

	// Update contact display name if available
	if senderContact.Profile.Name != "" {
		contact.DisplayName = senderContact.Profile.Name
		s.sm.DB.Save(contact)
	}

	// Save incoming message
	incomingMessage := &models.Message{
		MessageID:   message.ID,
		ContactID:   contact.ID,
		Content:     getMessageContent(message),
		MessageType: message.Type,
		Direction:   "incoming",
		Status:      "received",
		Timestamp:   time.Now(),
	}

	if err := s.sm.DB.Create(incomingMessage).Error; err != nil {
		return err
	}

	// Process auto-reply
	if err := s.sm.AutoReplyService.ProcessAutoReply(contact, incomingMessage); err != nil {
		logger.Log.WithError(err).Error("Failed to process auto-reply")
	}

	// Process custom commands
	if err := s.processCustomCommands(contact, incomingMessage); err != nil {
		logger.Log.WithError(err).Error("Failed to process custom commands")
	}

	// Process game commands
	if err := s.sm.GameService.ProcessGameCommand(contact, incomingMessage); err != nil {
		logger.Log.WithError(err).Error("Failed to process game command")
	}

	// Process utility commands
	if err := s.sm.UtilityService.ProcessUtilityCommand(contact, incomingMessage); err != nil {
		logger.Log.WithError(err).Error("Failed to process utility command")
	}

	// Process business commands
	if err := s.sm.BusinessService.ProcessBusinessCommand(contact, incomingMessage); err != nil {
		logger.Log.WithError(err).Error("Failed to process business command")
	}

	// Log analytics
	s.sm.AnalyticsService.LogEvent(contact.UserID, "message_received", 1, map[string]interface{}{
		"message_type": message.Type,
		"contact_id":   contact.ID,
	})

	return nil
}

func (s *WhatsAppService) processCustomCommands(contact *models.Contact, message *models.Message) error {
	// Check for custom commands
	commands, err := s.sm.DB.Where("user_id = ? AND is_active = ?", contact.UserID, true).Find(&models.CustomCommand{}).Error
	if err != nil {
		return err
	}

	for _, command := range commands {
		if shouldTriggerCommand(message.Content, command.Command, command.TriggerType) {
			// Send response
			_, err := s.SendMessage(contact.UserID, contact.PhoneNumber, command.Response, "text")
			return err
		}
	}

	return nil
}

func shouldTriggerCommand(message, command, triggerType string) bool {
	switch triggerType {
	case "exact":
		return strings.EqualFold(strings.TrimSpace(message), strings.TrimSpace(command))
	case "contains":
		return strings.Contains(strings.ToLower(message), strings.ToLower(command))
	case "regex":
		// Implement regex matching
		return false // Placeholder
	default:
		return false
	}
}

func getMessageContent(message *whatsapp.Message) string {
	switch message.Type {
	case "text":
		if message.Text != nil {
			return message.Text.Body
		}
	case "image":
		if message.Image != nil {
			return fmt.Sprintf("Image: %s", message.Image.Caption)
		}
	case "audio":
		if message.Audio != nil {
			return "Audio message"
		}
	case "video":
		if message.Video != nil {
			return "Video message"
		}
	case "document":
		if message.Document != nil {
			return "Document message"
		}
	}
	return ""
}

func (s *WhatsAppService) GetContacts(userID uuid.UUID) ([]models.Contact, error) {
	var contacts []models.Contact
	err := s.sm.DB.Where("user_id = ?", userID).Find(&contacts).Error
	return contacts, err
}

func (s *WhatsAppService) CreateGroup(userID uuid.UUID, name string, members []string) (*models.Group, error) {
	group := &models.Group{
		UserID:  userID,
		Name:    name,
		IsActive: true,
	}

	if err := s.sm.DB.Create(group).Error; err != nil {
		return nil, err
	}

	// Add members
	for _, member := range members {
		contact, err := s.sm.ContactService.FindOrCreateContact(member)
		if err != nil {
			continue
		}

		groupMember := &models.GroupMember{
			GroupID:   group.ID,
			ContactID: contact.ID,
			Role:      "member",
			JoinedAt:  time.Now(),
		}

		s.sm.DB.Create(groupMember)
	}

	return group, nil
}

func (s *WhatsAppService) GetGroups(userID uuid.UUID) ([]models.Group, error) {
	var groups []models.Group
	err := s.sm.DB.Where("user_id = ?", userID).Preload("Members").Find(&groups).Error
	return groups, err
}

func (s *WhatsAppService) UpdateMessageStatus(messageID string, status string) error {
	return s.sm.DB.Model(&models.Message{}).
		Where("message_id = ?", messageID).
		Update("status", status).Error
}

func (s *WhatsAppService) GetMessageStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total messages
	var totalMessages int64
	s.sm.DB.Model(&models.Message{}).Where("user_id = ?", userID).Count(&totalMessages)
	stats["total_messages"] = totalMessages

	// Messages by type
	var messagesByType []struct {
		MessageType string
		Count       int64
	}
	s.sm.DB.Model(&models.Message{}).
		Where("user_id = ?", userID).
		Select("message_type, count(*) as count").
		Group("message_type").
		Scan(&messagesByType)
	stats["messages_by_type"] = messagesByType

	// Messages by direction
	var messagesByDirection []struct {
		Direction string
		Count     int64
	}
	s.sm.DB.Model(&models.Message{}).
		Where("user_id = ?", userID).
		Select("direction, count(*) as count").
		Group("direction").
		Scan(&messagesByDirection)
	stats["messages_by_direction"] = messagesByDirection

	return stats, nil
}