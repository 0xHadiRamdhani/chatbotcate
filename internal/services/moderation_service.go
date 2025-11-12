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

type ModerationService struct {
	sm *ServiceManager
}

func (s *ModerationService) ProcessMessage(contact *models.Contact, message *models.Message) error {
	// Check for spam
	if s.isSpam(message.Content) {
		return s.handleSpam(contact, message)
	}

	// Check for blocked words
	if s.containsBlockedWords(message.Content, contact.UserID) {
		return s.handleBlockedWords(contact, message)
	}

	// Check for flood
	if s.isFlood(contact.ID) {
		return s.handleFlood(contact, message)
	}

	// Check for suspicious links
	if s.containsSuspiciousLinks(message.Content) {
		return s.handleSuspiciousLinks(contact, message)
	}

	// Check for inappropriate content
	if s.isInappropriate(message.Content) {
		return s.handleInappropriateContent(contact, message)
	}

	return nil
}

func (s *ModerationService) isSpam(content string) bool {
	// Check for spam patterns
	spamPatterns := []string{
		"http://", "https://", "www.",
		"buy now", "click here", "limited time",
		"free money", "make money fast", "work from home",
		"congratulations", "winner", "prize",
	}

	contentLower := strings.ToLower(content)
	for _, pattern := range spamPatterns {
		if strings.Contains(contentLower, pattern) {
			return true
		}
	}

	// Check for excessive caps
	if len(content) > 10 {
		upperCount := 0
		for _, char := range content {
			if char >= 'A' && char <= 'Z' {
				upperCount++
			}
		}
		if float64(upperCount)/float64(len(content)) > 0.7 {
			return true
		}
	}

	// Check for repeated characters
	repeatedCharPattern := regexp.MustCompile(`(.){5,}`)
	if repeatedCharPattern.MatchString(content) {
		return true
	}

	return false
}

func (s *ModerationService) containsBlockedWords(content string, userID uuid.UUID) bool {
	// Get blocked words for user
	var blockedWords []models.BlockedWord
	err := s.sm.DB.Where("user_id = ?", userID).Find(&blockedWords).Error
	if err != nil {
		return false
	}

	contentLower := strings.ToLower(content)
	for _, blockedWord := range blockedWords {
		if strings.Contains(contentLower, strings.ToLower(blockedWord.Word)) {
			return true
		}
	}

	return false
}

func (s *ModerationService) isFlood(contactID uuid.UUID) bool {
	// Check message frequency in last minute
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	
	var messageCount int64
	err := s.sm.DB.Model(&models.Message{}).
		Where("contact_id = ? AND direction = ? AND timestamp > ?", contactID, "incoming", oneMinuteAgo).
		Count(&messageCount).Error
	
	if err != nil {
		return false
	}

	// Consider it flood if more than 10 messages in 1 minute
	return messageCount > 10
}

func (s *ModerationService) containsSuspiciousLinks(content string) bool {
	// Check for suspicious link patterns
	suspiciousPatterns := []string{
		"bit.ly", "tinyurl.com", "goo.gl",
		"suspicious-domain.com", "phishing-site.com",
	}

	contentLower := strings.ToLower(content)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(contentLower, pattern) {
			return true
		}
	}

	return false
}

func (s *ModerationService) isInappropriate(content string) bool {
	// Basic inappropriate content detection
	inappropriateWords := []string{
		"spam", "scam", "fake", "fraud",
		"hack", "crack", "pirate",
	}

	contentLower := strings.ToLower(content)
	for _, word := range inappropriateWords {
		if strings.Contains(contentLower, word) {
			return true
		}
	}

	return false
}

func (s *ModerationService) handleSpam(contact *models.Contact, message *models.Message) error {
	// Block the message
	message.Status = "blocked"
	s.sm.DB.Save(message)

	// Send warning to user
	warningMessage := "⚠️ PERINGATAN ⚠️\n\nPesan Anda terdeteksi sebagai spam. Mohon kirim pesan yang relevan dan tidak mengandung promosi berlebihan."
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, warningMessage, false)
	if err != nil {
		return err
	}

	// Log incident
	s.logModerationIncident(contact.UserID, contact.ID, "spam_detected", message.Content)

	// Increment spam counter
	s.incrementSpamCounter(contact.ID)

	return nil
}

func (s *ModerationService) handleBlockedWords(contact *models.Contact, message *models.Message) error {
	// Block the message
	message.Status = "blocked"
	s.sm.DB.Save(message)

	// Send warning to user
	warningMessage := "⚠️ PERINGATAN ⚠️\n\nPesan Anda mengandung kata-kata yang tidak diizinkan. Mohon gunakan bahasa yang sopan dan tidak mengandung kata-kata sensitif."
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, warningMessage, false)
	if err != nil {
		return err
	}

	// Log incident
	s.logModerationIncident(contact.UserID, contact.ID, "blocked_words", message.Content)

	return nil
}

func (s *ModerationService) handleFlood(contact *models.Contact, message *models.Message) error {
	// Block the message
	message.Status = "blocked"
	s.sm.DB.Save(message)

	// Send warning to user
	warningMessage := "⚠️ PERINGATAN ⚠️\n\nAnda mengirim pesan terlalu cepat. Mohon tunggu beberapa saat sebelum mengirim pesan lagi."
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, warningMessage, false)
	if err != nil {
		return err
	}

	// Log incident
	s.logModerationIncident(contact.UserID, contact.ID, "flood_detected", message.Content)

	return nil
}

func (s *ModerationService) handleSuspiciousLinks(contact *models.Contact, message *models.Message) error {
	// Block the message
	message.Status = "blocked"
	s.sm.DB.Save(message)

	// Send warning to user
	warningMessage := "⚠️ PERINGATAN KEAMANAN ⚠️\n\nPesan Anda mengandung link yang mencurigakan. Untuk keamanan, link tersebut telah diblokir."
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, warningMessage, false)
	if err != nil {
		return err
	}

	// Log incident
	s.logModerationIncident(contact.UserID, contact.ID, "suspicious_link", message.Content)

	return nil
}

func (s *ModerationService) handleInappropriateContent(contact *models.Contact, message *models.Message) error {
	// Block the message
	message.Status = "blocked"
	s.sm.DB.Save(message)

	// Send warning to user
	warningMessage := "⚠️ PERINGATAN ⚠️\n\nPesan Anda mengandung konten yang tidak pantas. Mohon gunakan platform ini dengan bijak."
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, warningMessage, false)
	if err != nil {
		return err
	}

	// Log incident
	s.logModerationIncident(contact.UserID, contact.ID, "inappropriate_content", message.Content)

	return nil
}

func (s *ModerationService) logModerationIncident(userID, contactID uuid.UUID, incidentType, content string) {
	systemLog := &models.SystemLog{
		UserID:   userID,
		Level:    "warn",
		Message:  fmt.Sprintf("Moderation incident: %s", incidentType),
		Context:  fmt.Sprintf(`{"contact_id": "%s", "content": "%s"}`, contactID.String(), content),
	}

	s.sm.DB.Create(systemLog)
}

func (s *ModerationService) incrementSpamCounter(contactID uuid.UUID) {
	// Increment spam counter in Redis
	key := fmt.Sprintf("spam_count:%s", contactID.String())
	s.sm.Redis.Incr(s.sm.Redis.Context(), key)
	s.sm.Redis.Expire(s.sm.Redis.Context(), key, 24*time.Hour)
}

func (s *ModerationService) CreateBlockedWord(userID uuid.UUID, word, action string, replaceWith string) (*models.BlockedWord, error) {
	blockedWord := &models.BlockedWord{
		UserID:      userID,
		Word:        word,
		Action:      action,
		ReplaceWith: replaceWith,
	}

	if err := s.sm.DB.Create(blockedWord).Error; err != nil {
		return nil, err
	}

	return blockedWord, nil
}

func (s *ModerationService) GetBlockedWords(userID uuid.UUID) ([]models.BlockedWord, error) {
	var blockedWords []models.BlockedWord
	err := s.sm.DB.Where("user_id = ?", userID).Find(&blockedWords).Error
	return blockedWords, err
}

func (s *ModerationService) UpdateBlockedWord(id uuid.UUID, updates map[string]interface{}) error {
	return s.sm.DB.Model(&models.BlockedWord{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ModerationService) DeleteBlockedWord(id uuid.UUID) error {
	return s.sm.DB.Where("id = ?", id).Delete(&models.BlockedWord{}).Error
}

func (s *ModerationService) ReportSpam(reporterID, reportedID uuid.UUID, reason, evidence string) error {
	spamReport := &models.SpamReport{
		ReporterID: reporterID,
		ReportedID: reportedID,
		Reason:     reason,
		Evidence:   evidence,
		Status:     "pending",
	}

	if err := s.sm.DB.Create(spamReport).Error; err != nil {
		return err
	}

	// Notify admin
	s.notifyAdminSpamReport(spamReport)

	return nil
}

func (s *ModerationService) notifyAdminSpamReport(report *models.SpamReport) {
	// Get admin contacts
	var adminContacts []models.Contact
	s.sm.DB.Joins("JOIN users ON users.id = contacts.user_id").
		Where("users.is_admin = ?", true).
		Find(&adminContacts)

	reporter := &models.Contact{}
	s.sm.DB.Where("id = ?", report.ReporterID).First(reporter)

	reported := &models.Contact{}
	s.sm.DB.Where("id = ?", report.ReportedID).First(reported)

	message := fmt.Sprintf("🚨 LAPORAN SPAM BARU 🚨\n\nPelapor: %s\nDilaporkan: %s\nAlasan: %s\n\nSilakan tinjau laporan ini segera.", 
		reporter.DisplayName, reported.DisplayName, report.Reason)

	for _, admin := range adminContacts {
		s.sm.WhatsApp.SendTextMessage(admin.PhoneNumber, message, false)
	}
}

func (s *ModerationService) GetSpamReports(status string) ([]models.SpamReport, error) {
	var reports []models.SpamReport
	query := s.sm.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&reports).Error
	return reports, err
}

func (s *ModerationService) UpdateSpamReportStatus(id uuid.UUID, status string) error {
	return s.sm.DB.Model(&models.SpamReport{}).Where("id = ?", id).Update("status", status).Error
}

func (s *ModerationService) BlockContact(contactID uuid.UUID) error {
	return s.sm.DB.Model(&models.Contact{}).Where("id = ?", contactID).Update("is_blocked", true).Error
}

func (s *ModerationService) UnblockContact(contactID uuid.UUID) error {
	return s.sm.DB.Model(&models.Contact{}).Where("id = ?", contactID).Update("is_blocked", false).Error
}

func (s *ModerationService) GetBlockedContacts(userID uuid.UUID) ([]models.Contact, error) {
	var contacts []models.Contact
	err := s.sm.DB.Where("user_id = ? AND is_blocked = ?", userID, true).Find(&contacts).Error
	return contacts, err
}

func (s *ModerationService) SetMessageRateLimit(userID uuid.UUID, limit int) error {
	// Store rate limit in Redis
	key := fmt.Sprintf("rate_limit:%s", userID.String())
	return s.sm.Redis.Set(s.sm.Redis.Context(), key, limit, 0).Err()
}

func (s *ModerationService) GetMessageRateLimit(userID uuid.UUID) (int, error) {
	key := fmt.Sprintf("rate_limit:%s", userID.String())
	result, err := s.sm.Redis.Get(s.sm.Redis.Context(), key).Result()
	if err != nil {
		return 60, nil // Default 60 messages per minute
	}

	return strconv.Atoi(result)
}

func (s *ModerationService) CleanUpOldMessages(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	return s.sm.DB.Where("created_at < ?", cutoffDate).
		Delete(&models.Message{}).Error
}

func (s *ModerationService) GetModerationStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Blocked messages
	var blockedMessages int64
	s.sm.DB.Model(&models.Message{}).
		Where("user_id = ? AND status = ?", userID, "blocked").
		Count(&blockedMessages)
	stats["blocked_messages"] = blockedMessages

	// Blocked words count
	var blockedWords int64
	s.sm.DB.Model(&models.BlockedWord{}).
		Where("user_id = ?", userID).
		Count(&blockedWords)
	stats["blocked_words"] = blockedWords

	// Spam reports
	var spamReports int64
	s.sm.DB.Model(&models.SpamReport{}).
		Where("reporter_id IN (SELECT id FROM contacts WHERE user_id = ?)", userID).
		Count(&spamReports)
	stats["spam_reports"] = spamReports

	// Blocked contacts
	var blockedContacts int64
	s.sm.DB.Model(&models.Contact{}).
		Where("user_id = ? AND is_blocked = ?", userID, true).
		Count(&blockedContacts)
	stats["blocked_contacts"] = blockedContacts

	return stats, nil
}

func (s *ModerationService) CreateContentFilter(userID uuid.UUID, filterType, pattern string) error {
	// Create content filter based on type
	switch filterType {
	case "word":
		return s.CreateBlockedWord(userID, pattern, "block", "")
	case "regex":
		// Store regex pattern for content filtering
		return nil
	case "link":
		// Store link pattern for filtering
		return nil
	default:
		return fmt.Errorf("invalid filter type")
	}
}

func (s *ModerationService) EnableAutoModeration(userID uuid.UUID) error {
	// Enable auto-moderation for user
	return s.sm.DB.Model(&models.UserPreferences{}).
		Where("user_id = ?", userID).
		Update("enable_moderation", true).Error
}

func (s *ModerationService) DisableAutoModeration(userID uuid.UUID) error {
	// Disable auto-moderation for user
	return s.sm.DB.Model(&models.UserPreferences{}).
		Where("user_id = ?", userID).
		Update("enable_moderation", false).Error
}

func (s *ModerationService) ExportModerationSettings(userID uuid.UUID) ([]byte, error) {
	// Export moderation settings
	blockedWords, err := s.GetBlockedWords(userID)
	if err != nil {
		return nil, err
	}

	settings := map[string]interface{}{
		"blocked_words": blockedWords,
		"rate_limit":    s.GetMessageRateLimit(userID),
	}

	return json.Marshal(settings)
}

func (s *ModerationService) ImportModerationSettings(userID uuid.UUID, data []byte) error {
	// Import moderation settings
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	// Process imported settings
	if blockedWords, ok := settings["blocked_words"].([]interface{}); ok {
		for _, word := range blockedWords {
			if wordMap, ok := word.(map[string]interface{}); ok {
				wordText := wordMap["word"].(string)
				action := wordMap["action"].(string)
				s.CreateBlockedWord(userID, wordText, action, "")
			}
		}
	}

	return nil
}