package services

import (
	"fmt"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	sm *ServiceManager
}

func (s *UserService) CreateUser(username, email, password, phoneNumber string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:    username,
		Email:       email,
		Password:    string(hashedPassword),
		PhoneNumber: phoneNumber,
		DisplayName: username,
		IsActive:    true,
		IsAdmin:     false,
		Points:      0,
		Level:       1,
	}

	if err := s.sm.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Create user preferences
	preferences := &models.UserPreferences{
		UserID:           user.ID,
		Language:         "id",
		Timezone:         "Asia/Jakarta",
		EnableGames:      true,
		EnableBusiness:   true,
		EnableUtils:      true,
		EnableModeration: true,
	}

	if err := s.sm.DB.Create(preferences).Error; err != nil {
		return nil, err
	}

	// Log analytics
	s.sm.AnalyticsService.LogEvent(user.ID, "user_created", 1, map[string]interface{}{
		"username": username,
		"email":    email,
	})

	return user, nil
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	err := s.sm.DB.Where("username = ? OR email = ?", username, username).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Update last login
	user.LastLoginAt = &time.Now()
	s.sm.DB.Save(&user)

	// Log analytics
	s.sm.AnalyticsService.LogEvent(user.ID, "user_login", 1, map[string]interface{}{
		"username": username,
	})

	return &user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.sm.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.sm.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, updates map[string]interface{}) error {
	return s.sm.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserService) UpdateUserPreferences(userID uuid.UUID, preferences map[string]interface{}) error {
	return s.sm.DB.Model(&models.UserPreferences{}).Where("user_id = ?", userID).Updates(preferences).Error
}

func (s *UserService) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	var preferences models.UserPreferences
	err := s.sm.DB.Where("user_id = ?", userID).First(&preferences).Error
	if err != nil {
		return nil, err
	}

	return &preferences, nil
}

func (s *UserService) AddPoints(userID uuid.UUID, points int) error {
	user := &models.User{}
	if err := s.sm.DB.Where("id = ?", userID).First(user).Error; err != nil {
		return err
	}

	user.Points += points

	// Check level up
	newLevel := s.calculateLevel(user.Points)
	if newLevel > user.Level {
		user.Level = newLevel
		// Send level up notification
		s.sendLevelUpNotification(user, newLevel)
	}

	return s.sm.DB.Save(user).Error
}

func (s *UserService) calculateLevel(points int) int {
	// Simple level calculation
	if points >= 1000 {
		return 5
	} else if points >= 500 {
		return 4
	} else if points >= 200 {
		return 3
	} else if points >= 50 {
		return 2
	}
	return 1
}

func (s *UserService) sendLevelUpNotification(user *models.User, newLevel int) {
	// Get user's primary contact
	contact := &models.Contact{}
	err := s.sm.DB.Where("user_id = ?", user.ID).Order("created_at ASC").First(contact).Error
	if err != nil {
		return
	}

	message := fmt.Sprintf("🎉 LEVEL UP! 🎉\n\nSelamat! Anda naik ke level %d!\nPoin Anda: %d\n\nTerus gunakan bot kami untuk mendapatkan lebih banyak poin dan keuntungan!", newLevel, user.Points)
	_, err = s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to send level up notification")
	}
}

func (s *UserService) GetUserStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get user
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	stats["points"] = user.Points
	stats["level"] = user.Level
	stats["last_login"] = user.LastLoginAt

	// Total messages
	var totalMessages int64
	s.sm.DB.Model(&models.Message{}).Where("user_id = ?", userID).Count(&totalMessages)
	stats["total_messages"] = totalMessages

	// Total contacts
	var totalContacts int64
	s.sm.DB.Model(&models.Contact{}).Where("user_id = ?", userID).Count(&totalContacts)
	stats["total_contacts"] = totalContacts

	// Total orders (if business feature enabled)
	var totalOrders int64
	s.sm.DB.Model(&models.Order{}).Where("user_id = ?", userID).Count(&totalOrders)
	stats["total_orders"] = totalOrders

	return stats, nil
}

func (s *UserService) GetUsers(limit int, offset int) ([]models.User, error) {
	var users []models.User
	query := s.sm.DB

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&users).Error
	return users, err
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	// Soft delete user
	return s.sm.DB.Where("id = ?", id).Delete(&models.User{}).Error
}

func (s *UserService) DeactivateUser(id uuid.UUID) error {
	return s.sm.DB.Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}

func (s *UserService) ActivateUser(id uuid.UUID) error {
	return s.sm.DB.Model(&models.User{}).Where("id = ?", id).Update("is_active", true).Error
}

func (s *UserService) MakeAdmin(id uuid.UUID) error {
	return s.sm.DB.Model(&models.User{}).Where("id = ?", id).Update("is_admin", true).Error
}

func (s *UserService) RemoveAdmin(id uuid.UUID) error {
	return s.sm.DB.Model(&models.User{}).Where("id = ?", id).Update("is_admin", false).Error
}

func (s *UserService) GetUserActivity(userID uuid.UUID, days int) ([]models.SystemLog, error) {
	var logs []models.SystemLog
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	err := s.sm.DB.Where("user_id = ? AND created_at > ?", userID, cutoffDate).
		Order("created_at DESC").
		Find(&logs).Error
	
	return logs, err
}

func (s *UserService) ExportUserData(userID uuid.UUID) ([]byte, error) {
	// Export user data including preferences, contacts, messages, etc.
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	preferences, err := s.GetUserPreferences(userID)
	if err != nil {
		return nil, err
	}

	// Get contacts
	var contacts []models.Contact
	s.sm.DB.Where("user_id = ?", userID).Find(&contacts)

	// Get messages (last 1000)
	var messages []models.Message
	s.sm.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(1000).Find(&messages)

	data := map[string]interface{}{
		"user":        user,
		"preferences": preferences,
		"contacts":    contacts,
		"messages":    messages,
	}

	return json.Marshal(data)
}

func (s *UserService) ImportUserData(userID uuid.UUID, data []byte) error {
	// Import user data
	var importData map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		return err
	}

	// Process imported data
	if contacts, ok := importData["contacts"].([]interface{}); ok {
		for _, contact := range contacts {
			if contactMap, ok := contact.(map[string]interface{}); ok {
				// Create contact
				newContact := &models.Contact{
					UserID:      userID,
					PhoneNumber: contactMap["phone_number"].(string),
					DisplayName: contactMap["display_name"].(string),
				}
				s.sm.DB.Create(newContact)
			}
		}
	}

	return nil
}

func (s *UserService) CreateUserFromWhatsApp(phoneNumber, displayName string) (*models.User, error) {
	// Create user from WhatsApp contact
	username := fmt.Sprintf("user_%s", phoneNumber)
	email := fmt.Sprintf("%s@whatsapp.bot", phoneNumber)
	password := generateRandomPassword()

	user, err := s.CreateUser(username, email, password, phoneNumber)
	if err != nil {
		return nil, err
	}

	// Set display name
	user.DisplayName = displayName
	s.sm.DB.Save(user)

	return user, nil
}

func (s *UserService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := s.sm.DB.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateLastActivity(userID uuid.UUID) error {
	return s.sm.DB.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", time.Now()).Error
}

func (s *UserService) GetUserDashboard(userID uuid.UUID) (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})

	// User info
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	dashboard["user"] = user

	// Recent messages
	var recentMessages []models.Message
	s.sm.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentMessages)
	dashboard["recent_messages"] = recentMessages

	// Active reminders
	var activeReminders []models.Reminder
	s.sm.DB.Where("user_id = ? AND status = ?", userID, "active").
		Order("remind_at ASC").
		Limit(3).
		Find(&activeReminders)
	dashboard["active_reminders"] = activeReminders

	// Recent orders (if business enabled)
	var recentOrders []models.Order
	s.sm.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(3).
		Find(&recentOrders)
	dashboard["recent_orders"] = recentOrders

	// Game leaderboard position
	var gameScores []models.GameScore
	s.sm.DB.Where("user_id = ?", userID).
		Order("score DESC").
		Find(&gameScores)
	dashboard["game_scores"] = gameScores

	return dashboard, nil
}

func generateRandomPassword() string {
	// Generate random password
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	password := make([]byte, 12)
	
	for i := 0; i < 12; i++ {
		password[i] = chars[rand.Intn(len(chars))]
	}
	
	return string(password)
}