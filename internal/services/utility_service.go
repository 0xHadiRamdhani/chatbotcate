package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/pkg/logger"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type UtilityService struct {
	db *Database
}

func NewUtilityService(db *Database) *UtilityService {
	return &UtilityService{
		db: db,
	}
}

// CreateQRCode creates a QR code
func (s *UtilityService) CreateQRCode(data string, size int, format string) (string, error) {
	var png []byte
	png, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		logger.Error("Failed to generate QR code", err)
		return "", err
	}

	if format == "png" {
		return string(png), nil
	} else if format == "svg" {
		// Convert PNG to SVG (simplified implementation)
		svg := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
			<image href="data:image/png;base64,%s" width="%d" height="%d"/>
		</svg>`, size, size, utils.GenerateRandomHex(len(png)), size, size)
		return svg, nil
	}

	return string(png), nil
}

// CreateShortLink creates a short link
func (s *UtilityService) CreateShortLink(originalURL string, customAlias string, expiryDays int) (*models.ShortLink, error) {
	// Generate short alias if not provided
	alias := customAlias
	if alias == "" {
		alias = utils.GenerateRandomString(6)
	}

	// Check if alias already exists
	var existing models.ShortLink
	if err := s.db.DB.Where("alias = ?", alias).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("alias already exists")
	}

	shortLink := &models.ShortLink{
		ID:          uuid.New(),
		OriginalURL: originalURL,
		Alias:       alias,
		ExpiryAt:    time.Now().AddDate(0, 0, expiryDays),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.DB.Create(shortLink).Error; err != nil {
		logger.Error("Failed to create short link", err)
		return nil, err
	}

	return shortLink, nil
}

// GetShortLink gets a short link by alias
func (s *UtilityService) GetShortLink(alias string) (*models.ShortLink, error) {
	var shortLink models.ShortLink
	if err := s.db.DB.Where("alias = ?", alias).First(&shortLink).Error; err != nil {
		logger.Error("Failed to get short link", err)
		return nil, err
	}

	// Check if expired
	if shortLink.ExpiryAt.Before(time.Now()) {
		return nil, fmt.Errorf("short link expired")
	}

	// Update click count
	shortLink.ClickCount++
	shortLink.LastClickedAt = time.Now()
	if err := s.db.DB.Save(&shortLink).Error; err != nil {
		logger.Error("Failed to update short link clicks", err)
	}

	return &shortLink, nil
}

// ConvertCurrency converts currency using external API
func (s *UtilityService) ConvertCurrency(amount float64, from string, to string) (float64, error) {
	// In a real implementation, you would use an external API like ExchangeRate-API
	// For demo purposes, we'll use a mock conversion rate
	apiURL := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from)

	resp, err := http.Get(apiURL)
	if err != nil {
		logger.Error("Failed to fetch exchange rates", err)
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error("Failed to decode exchange rates", err)
		return 0, err
	}

	rates := data["rates"].(map[string]interface{})
	rate := rates[to].(float64)

	return amount * rate, nil
}

// GetWeather gets weather information
func (s *UtilityService) GetWeather(city string, country string) (map[string]interface{}, error) {
	// In a real implementation, you would use an external weather API
	// For demo purposes, we'll return mock data
	weather := map[string]interface{}{
		"city":        city,
		"country":     country,
		"temperature": 25,
		"condition":   "Sunny",
		"humidity":    60,
		"wind_speed":  10,
		"forecast": []map[string]interface{}{
			{"day": "Today", "high": 28, "low": 20, "condition": "Sunny"},
			{"day": "Tomorrow", "high": 26, "low": 18, "condition": "Cloudy"},
			{"day": "Day 3", "high": 24, "low": 16, "condition": "Rainy"},
		},
	}

	return weather, nil
}

// TranslateText translates text using external API
func (s *UtilityService) TranslateText(text string, targetLang string, sourceLang string) (string, error) {
	// In a real implementation, you would use Google Translate API or similar
	// For demo purposes, we'll return a mock translation
	translations := map[string]map[string]string{
		"en": {
			"hello": "hello",
			"goodbye": "goodbye",
			"thank you": "thank you",
		},
		"id": {
			"hello": "halo",
			"goodbye": "selamat tinggal",
			"thank you": "terima kasih",
		},
		"es": {
			"hello": "hola",
			"goodbye": "adiós",
			"thank you": "gracias",
		},
	}

	// Simple mock translation
	if sourceLang == "" {
		sourceLang = "en"
	}

	if trans, ok := translations[targetLang]; ok {
		if translated, ok := trans[text]; ok {
			return translated, nil
		}
	}

	// Return original text if no translation found
	return text, nil
}

// GetLocationInfo gets location information
func (s *UtilityService) GetLocationInfo(latitude float64, longitude float64) (map[string]interface{}, error) {
	// In a real implementation, you would use Google Maps API or similar
	// For demo purposes, we'll return mock data
	location := map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
		"address":   "Jl. Sudirman No. 123, Jakarta",
		"city":      "Jakarta",
		"country":   "Indonesia",
		"timezone":  "Asia/Jakarta",
		"postal_code": "12345",
	}

	return location, nil
}

// CreatePoll creates a poll
func (s *UtilityService) CreatePoll(question string, options []string, userID uuid.UUID) (*models.Poll, error) {
	poll := &models.Poll{
		ID:        uuid.New(),
		Question:  question,
		Options:   options,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.DB.Create(poll).Error; err != nil {
		logger.Error("Failed to create poll", err)
		return nil, err
	}

	return poll, nil
}

// VotePoll votes on a poll
func (s *UtilityService) VotePoll(pollID uuid.UUID, optionID int, userID uuid.UUID) error {
	// Check if user already voted
	var existingVote models.PollVote
	if err := s.db.DB.Where("poll_id = ? AND user_id = ?", pollID, userID).First(&existingVote).Error; err == nil {
		return fmt.Errorf("user already voted on this poll")
	}

	vote := &models.PollVote{
		ID:        uuid.New(),
		PollID:    pollID,
		OptionID:  optionID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := s.db.DB.Create(vote).Error; err != nil {
		logger.Error("Failed to create poll vote", err)
		return err
	}

	return nil
}

// GetPollResults gets poll results
func (s *UtilityService) GetPollResults(pollID uuid.UUID) (map[string]interface{}, error) {
	var poll models.Poll
	if err := s.db.DB.Where("id = ?", pollID).First(&poll).Error; err != nil {
		logger.Error("Failed to get poll", err)
		return nil, err
	}

	// Get vote counts for each option
	results := make([]map[string]interface{}, len(poll.Options))
	for i, option := range poll.Options {
		var count int64
		s.db.DB.Model(&models.PollVote{}).Where("poll_id = ? AND option_id = ?", pollID, i).Count(&count)
		
		results[i] = map[string]interface{}{
			"option": option,
			"votes":  count,
		}
	}

	// Get total votes
	var totalVotes int64
	s.db.DB.Model(&models.PollVote{}).Where("poll_id = ?", pollID).Count(&totalVotes)

	return map[string]interface{}{
		"poll_id":     pollID,
		"question":    poll.Question,
		"options":     results,
		"total_votes": totalVotes,
		"created_at":  poll.CreatedAt,
	}, nil
}

// CreateReminder creates a reminder
func (s *UtilityService) CreateReminder(title string, description string, remindAt time.Time, userID uuid.UUID, repeat string) (*models.Reminder, error) {
	reminder := &models.Reminder{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		RemindAt:    remindAt,
		UserID:      userID,
		Repeat:      repeat,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.DB.Create(reminder).Error; err != nil {
		logger.Error("Failed to create reminder", err)
		return nil, err
	}

	return reminder, nil
}

// GetReminders gets user reminders
func (s *UtilityService) GetReminders(userID uuid.UUID) ([]models.Reminder, error) {
	var reminders []models.Reminder
	if err := s.db.DB.Where("user_id = ? AND is_active = ?", userID, true).Order("remind_at asc").Find(&reminders).Error; err != nil {
		logger.Error("Failed to get reminders", err)
		return nil, err
	}

	return reminders, nil
}

// DeleteReminder deletes a reminder
func (s *UtilityService) DeleteReminder(reminderID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", reminderID).Delete(&models.Reminder{}).Error; err != nil {
		logger.Error("Failed to delete reminder", err)
		return err
	}

	return nil
}

// CreateNote creates a note
func (s *UtilityService) CreateNote(title string, content string, userID uuid.UUID, category string, tags []string) (*models.Note, error) {
	note := &models.Note{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		UserID:    userID,
		Category:  category,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.DB.Create(note).Error; err != nil {
		logger.Error("Failed to create note", err)
		return nil, err
	}

	return note, nil
}

// GetNotes gets user notes
func (s *UtilityService) GetNotes(userID uuid.UUID, category string, tag string) ([]models.Note, error) {
	query := s.db.DB.Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	var notes []models.Note
	if err := query.Order("created_at desc").Find(&notes).Error; err != nil {
		logger.Error("Failed to get notes", err)
		return nil, err
	}

	return notes, nil
}

// UpdateNote updates a note
func (s *UtilityService) UpdateNote(noteID uuid.UUID, title string, content string, category string, tags []string) (*models.Note, error) {
	var note models.Note
	if err := s.db.DB.Where("id = ?", noteID).First(&note).Error; err != nil {
		logger.Error("Failed to get note for update", err)
		return nil, err
	}

	// Update fields
	if title != "" {
		note.Title = title
	}
	if content != "" {
		note.Content = content
	}
	if category != "" {
		note.Category = category
	}
	if tags != nil {
		note.Tags = tags
	}
	note.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&note).Error; err != nil {
		logger.Error("Failed to update note", err)
		return nil, err
	}

	return &note, nil
}

// DeleteNote deletes a note
func (s *UtilityService) DeleteNote(noteID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", noteID).Delete(&models.Note{}).Error; err != nil {
		logger.Error("Failed to delete note", err)
		return err
	}

	return nil
}

// SearchNotes searches notes
func (s *UtilityService) SearchNotes(userID uuid.UUID, query string) ([]models.Note, error) {
	var notes []models.Note
	searchQuery := "%" + query + "%"
	
	if err := s.db.DB.Where("user_id = ? AND (title LIKE ? OR content LIKE ? OR tags LIKE ?)", 
		userID, searchQuery, searchQuery, searchQuery).Order("created_at desc").Find(&notes).Error; err != nil {
		logger.Error("Failed to search notes", err)
		return nil, err
	}

	return notes, nil
}

// CreateTimer creates a timer
func (s *UtilityService) CreateTimer(name string, duration int, userID uuid.UUID) (*models.Timer, error) {
	timer := &models.Timer{
		ID:        uuid.New(),
		Name:      name,
		Duration:  duration,
		UserID:    userID,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Duration(duration) * time.Second),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.DB.Create(timer).Error; err != nil {
		logger.Error("Failed to create timer", err)
		return nil, err
	}

	return timer, nil
}

// GetTimer gets timer status
func (s *UtilityService) GetTimer(timerID uuid.UUID) (*models.Timer, error) {
	var timer models.Timer
	if err := s.db.DB.Where("id = ?", timerID).First(&timer).Error; err != nil {
		logger.Error("Failed to get timer", err)
		return nil, err
	}

	// Update remaining time
	if timer.IsActive && timer.EndTime.After(time.Now()) {
		remaining := int(timer.EndTime.Sub(time.Now()).Seconds())
		if remaining < 0 {
			remaining = 0
			timer.IsActive = false
			s.db.DB.Save(&timer)
		}
	}

	return &timer, nil
}

// StopTimer stops a timer
func (s *UtilityService) StopTimer(timerID uuid.UUID) error {
	var timer models.Timer
	if err := s.db.DB.Where("id = ?", timerID).First(&timer).Error; err != nil {
		logger.Error("Failed to get timer for stop", err)
		return err
	}

	timer.IsActive = false
	timer.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&timer).Error; err != nil {
		logger.Error("Failed to stop timer", err)
		return err
	}

	return nil
}

// CreateFileUpload creates file upload metadata
func (s *UtilityService) CreateFileUpload(filename string, fileSize int64, fileType string, userID uuid.UUID) (*models.FileUpload, error) {
	upload := &models.FileUpload{
		ID:         uuid.New(),
		Filename:   filename,
		FileSize:   fileSize,
		FileType:   fileType,
		UserID:     userID,
		UploadURL:  fmt.Sprintf("/uploads/%s/%s", userID.String(), filename),
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.DB.Create(upload).Error; err != nil {
		logger.Error("Failed to create file upload", err)
		return nil, err
	}

	return upload, nil
}

// GetFileUpload gets file upload info
func (s *UtilityService) GetFileUpload(uploadID uuid.UUID) (*models.FileUpload, error) {
	var upload models.FileUpload
	if err := s.db.DB.Where("id = ?", uploadID).First(&upload).Error; err != nil {
		logger.Error("Failed to get file upload", err)
		return nil, err
	}

	return &upload, nil
}

// DeleteFileUpload deletes file upload
func (s *UtilityService) DeleteFileUpload(uploadID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", uploadID).Delete(&models.FileUpload{}).Error; err != nil {
		logger.Error("Failed to delete file upload", err)
		return err
	}

	return nil
}