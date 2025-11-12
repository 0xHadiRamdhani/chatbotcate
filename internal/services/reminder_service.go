package services

import (
	"fmt"
	"strings"
	"time"

	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/logger"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

type ReminderService struct {
	sm *ServiceManager
}

func (s *ReminderService) ProcessReminders() error {
	// Get active reminders that are due
	now := time.Now()
	var reminders []models.Reminder
	
	err := s.sm.DB.Where("status = ? AND remind_at <= ?", "active", now).Find(&reminders).Error
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		if err := s.sendReminder(&reminder); err != nil {
			logger.Log.WithError(err).Error("Failed to send reminder")
			continue
		}

		// Update reminder status
		if reminder.IsRecurring {
			// Schedule next reminder
			nextRemindAt := s.calculateNextReminder(reminder)
			reminder.RemindAt = nextRemindAt
		} else {
			reminder.Status = "completed"
			reminder.CompletedAt = &now
		}

		s.sm.DB.Save(&reminder)
	}

	return nil
}

func (s *ReminderService) sendReminder(reminder *models.Reminder) error {
	// Get contact information
	contact := &models.Contact{}
	if err := s.sm.DB.Where("id = ?", reminder.ContactID).First(contact).Error; err != nil {
		return err
	}

	// Format reminder message
	message := fmt.Sprintf("🔔 PENGINGAT 🔔\n\n%s\n\n%s", reminder.Title, reminder.Description)
	
	if reminder.IsRecurring {
		message += fmt.Sprintf("\n\n⏰ Pengingat ini akan diulang: %s", reminder.RecurringType)
	}

	// Send reminder
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	if err != nil {
		return err
	}

	// Log analytics
	s.sm.AnalyticsService.LogEvent(reminder.UserID, "reminder_sent", 1, map[string]interface{}{
		"reminder_id": reminder.ID,
		"recurring":   reminder.IsRecurring,
	})

	return nil
}

func (s *ReminderService) calculateNextReminder(reminder models.Reminder) time.Time {
	now := time.Now()
	
	switch reminder.RecurringType {
	case "daily":
		return now.Add(24 * time.Hour)
	case "weekly":
		return now.Add(7 * 24 * time.Hour)
	case "monthly":
		return now.AddDate(0, 1, 0)
	default:
		return now.Add(24 * time.Hour)
	}
}

func (s *ReminderService) CreateReminder(userID, contactID uuid.UUID, title, description string, remindAt time.Time, isRecurring bool, recurringType string) (*models.Reminder, error) {
	reminder := &models.Reminder{
		UserID:        userID,
		ContactID:     contactID,
		Title:         title,
		Description:   description,
		RemindAt:      remindAt,
		IsRecurring:   isRecurring,
		RecurringType: recurringType,
		Status:        "active",
	}

	if err := s.sm.DB.Create(reminder).Error; err != nil {
		return nil, err
	}

	return reminder, nil
}

func (s *ReminderService) GetReminders(userID uuid.UUID, status string) ([]models.Reminder, error) {
	var reminders []models.Reminder
	query := s.sm.DB.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("remind_at ASC")
	err := query.Find(&reminders).Error
	return reminders, err
}

func (s *ReminderService) UpdateReminder(id uuid.UUID, updates map[string]interface{}) error {
	return s.sm.DB.Model(&models.Reminder{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ReminderService) DeleteReminder(id uuid.UUID) error {
	return s.sm.DB.Where("id = ?", id).Delete(&models.Reminder{}).Error
}

func (s *ReminderService) CancelReminder(id uuid.UUID) error {
	return s.sm.DB.Model(&models.Reminder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "cancelled",
		"completed_at": time.Now(),
	}).Error
}

func (s *ReminderService) ProcessReminderCommand(contact *models.Contact, content string) error {
	// Parse reminder command
	words := strings.Fields(content)
	var action string
	var reminderText string
	var timeStr string

	// Find action and parameters
	for i, word := range words {
		if word == "reminder" || word == "pengingat" {
			if i+1 < len(words) {
				action = words[i+1]
				if i+2 < len(words) {
					reminderText = strings.Join(words[i+2:], " ")
				}
			}
			break
		}
	}

	switch action {
	case "buat", "create":
		return s.handleCreateReminder(contact, reminderText)
	case "list", "daftar":
		return s.handleListReminders(contact)
	case "hapus", "delete":
		return s.handleDeleteReminder(contact, reminderText)
	default:
		return s.handleReminderHelp(contact)
	}
}

func (s *ReminderService) handleCreateReminder(contact *models.Contact, text string) error {
	// Parse reminder text: "besok 08:00 meeting dengan client"
	words := strings.Fields(text)
	var reminderTime time.Time
	var title string
	var isRecurring bool
	var recurringType string

	// Parse time and date
	for i, word := range words {
		switch word {
		case "besok":
			reminderTime = time.Now().Add(24 * time.Hour)
			if i+1 < len(words) {
				// Parse time
				if t, err := time.Parse("HH:mm", words[i+1]); err == nil {
					reminderTime = time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), 
						t.Hour(), t.Minute(), 0, 0, reminderTime.Location())
				}
			}
		case "mingguan":
			isRecurring = true
			recurringType = "weekly"
		case "harian":
			isRecurring = true
			recurringType = "daily"
		case "bulanan":
			isRecurring = true
			recurringType = "monthly"
		}
	}

	// Extract title
	if reminderTime.IsZero() {
		// Default to tomorrow at 9 AM
		reminderTime = time.Now().Add(24 * time.Hour)
		reminderTime = time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), 9, 0, 0, 0, reminderTime.Location())
		title = text
	} else {
		// Extract remaining text as title
		title = strings.Join(words, " ")
	}

	// Create reminder
	reminder, err := s.CreateReminder(contact.UserID, contact.ID, "Pengingat: "+title, title, reminderTime, isRecurring, recurringType)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("✅ PENGINGAT DIBUAT ✅\n\nJudul: %s\nWaktu: %s\n\nPengingat akan dikirim pada waktu yang ditentukan.", title, reminderTime.Format("dd-MM-yyyy HH:mm"))
	if isRecurring {
		message += fmt.Sprintf("\nPengingat akan diulang: %s", recurringType)
	}

	_, err = s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	return err
}

func (s *ReminderService) handleListReminders(contact *models.Contact) error {
	reminders, err := s.GetReminders(contact.UserID, "active")
	if err != nil {
		return err
	}

	if len(reminders) == 0 {
		message := "📋 DAFTAR PENGINGAT 📋\n\nAnda belum memiliki pengingat aktif.\n\nUntuk membuat pengingat, ketik: 'reminder buat besok 08:00 meeting'"
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
		return err
	}

	message := "📋 DAFTAR PENGINGAT ANDA 📋\n\n"
	for i, reminder := range reminders {
		message += fmt.Sprintf("%d. %s\n", i+1, reminder.Title)
		message += fmt.Sprintf("   Waktu: %s\n", reminder.RemindAt.Format("dd-MM-yyyy HH:mm"))
		if reminder.IsRecurring {
			message += fmt.Sprintf("   Diulang: %s\n", reminder.RecurringType)
		}
		message += "\n"
	}

	message += "Untuk menghapus pengingat, ketik: 'reminder hapus [nomor]'"

	_, err = s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	return err
}

func (s *ReminderService) handleDeleteReminder(contact *models.Contact, text string) error {
	// Parse reminder number
	words := strings.Fields(text)
	var reminderNumber int

	for _, word := range words {
		if num, err := strconv.Atoi(word); err == nil {
			reminderNumber = num
			break
		}
	}

	if reminderNumber == 0 {
		message := "❌ FORMAT SALAH ❌\n\nGunakan: 'reminder hapus [nomor_pengingat]'\n\nContoh: 'reminder hapus 1'"
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
		return err
	}

	// Get reminders
	reminders, err := s.GetReminders(contact.UserID, "active")
	if err != nil {
		return err
	}

	if reminderNumber > len(reminders) {
		message := "❌ PENGINGAT TIDAK DITEMUKAN ❌\n\nNomor pengingat tidak valid."
		_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
		return err
	}

	reminder := reminders[reminderNumber-1]
	if err := s.DeleteReminder(reminder.ID); err != nil {
		return err
	}

	message := fmt.Sprintf("✅ PENGINGAT DIHAPUS ✅\n\nPengingat '%s' telah dihapus.", reminder.Title)
	_, err = s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	return err
}

func (s *ReminderService) handleReminderHelp(contact *models.Contact) error {
	message := "🔔 BANTUAN PENGINGAT 🔔\n\nPerintah yang tersedia:\n\n• 'reminder buat [teks]' - Buat pengingat\n• 'reminder list' - Lihat daftar pengingat\n• 'reminder hapus [nomor]' - Hapus pengingat\n\nContoh:\n• 'reminder buat besok 08:00 meeting dengan client'\n• 'reminder buat harian minum vitamin'\n• 'reminder list'\n• 'reminder hapus 1'"
	
	_, err := s.sm.WhatsApp.SendTextMessage(contact.PhoneNumber, message, false)
	return err
}

func (s *ReminderService) CreateToDoList(userID uuid.UUID, title string, items []string) error {
	// Create to-do list as multiple reminders
	for i, item := range items {
		reminder := &models.Reminder{
			UserID:      userID,
			Title:       fmt.Sprintf("%s - Item %d", title, i+1),
			Description: item,
			RemindAt:    time.Now().Add(24 * time.Hour), // Default reminder time
			Status:      "active",
		}

		if err := s.sm.DB.Create(reminder).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *ReminderService) ScheduleMessage(userID, contactID uuid.UUID, content string, scheduledTime time.Time) error {
	// Create reminder for scheduled message
	reminder := &models.Reminder{
		UserID:      userID,
		ContactID:   contactID,
		Title:       "Pesan Terjadwal",
		Description: content,
		RemindAt:    scheduledTime,
		Status:      "active",
	}

	if err := s.sm.DB.Create(reminder).Error; err != nil {
		return err
	}

	return nil
}

func (s *ReminderService) GetReminderStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total reminders
	var total int64
	s.sm.DB.Model(&models.Reminder{}).Where("user_id = ?", userID).Count(&total)
	stats["total_reminders"] = total

	// Active reminders
	var active int64
	s.sm.DB.Model(&models.Reminder{}).Where("user_id = ? AND status = ?", userID, "active").Count(&active)
	stats["active_reminders"] = active

	// Completed reminders
	var completed int64
	s.sm.DB.Model(&models.Reminder{}).Where("user_id = ? AND status = ?", userID, "completed").Count(&completed)
	stats["completed_reminders"] = completed

	// Recurring reminders
	var recurring int64
	s.sm.DB.Model(&models.Reminder{}).Where("user_id = ? AND is_recurring = ?", userID, true).Count(&recurring)
	stats["recurring_reminders"] = recurring

	return stats, nil
}

func (s *ReminderService) BulkCreateReminders(userID uuid.UUID, reminders []models.Reminder) error {
	tx := s.sm.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, reminder := range reminders {
		reminder.UserID = userID
		if err := tx.Create(&reminder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *ReminderService) ExportReminders(userID uuid.UUID) ([]byte, error) {
	reminders, err := s.GetReminders(userID, "")
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(reminders)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ReminderService) ImportReminders(userID uuid.UUID, data []byte) error {
	var reminders []models.Reminder
	if err := json.Unmarshal(data, &reminders); err != nil {
		return err
	}

	return s.BulkCreateReminders(userID, reminders)
}

func (s *ReminderService) SetUpDailyReminders(userID uuid.UUID, reminders []DailyReminderConfig) error {
	for _, config := range reminders {
		reminderTime, _ := time.Parse("HH:mm", config.Time)
		reminderTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 
			reminderTime.Hour(), reminderTime.Minute(), 0, 0, time.Now().Location())

		reminder := &models.Reminder{
			UserID:        userID,
			Title:         config.Title,
			Description:   config.Description,
			RemindAt:      reminderTime,
			IsRecurring:   true,
			RecurringType: "daily",
			Status:        "active",
		}

		if err := s.sm.DB.Create(reminder).Error; err != nil {
			return err
		}
	}

	return nil
}

type DailyReminderConfig struct {
	Title       string
	Description string
	Time        string
}

func (s *ReminderService) CleanUpOldReminders(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	return s.sm.DB.Where("status = ? AND completed_at < ?", "completed", cutoffDate).
		Delete(&models.Reminder{}).Error
}