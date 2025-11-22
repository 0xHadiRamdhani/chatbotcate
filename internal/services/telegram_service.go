package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/pkg/logger"
	"kilocode.dev/whatsapp-bot/pkg/telegram"
)

type TelegramService struct {
	db     *Database
	client *telegram.Client
}

func NewTelegramService(db *Database, apiKey string) *TelegramService {
	return &TelegramService{
		db:     db,
		client: telegram.NewClient(apiKey),
	}
}

// SendMessage sends a text message to Telegram
func (s *TelegramService) SendMessage(chatID int64, text string) error {
	return s.client.SendMessage(chatID, text, "")
}

// SendMessageWithMarkup sends a message with custom keyboard markup
func (s *TelegramService) SendMessageWithMarkup(chatID int64, text string, markup interface{}) error {
	return s.client.SendMessageWithMarkup(chatID, text, markup)
}

// SendPhoto sends a photo to Telegram
func (s *TelegramService) SendPhoto(chatID int64, photoURL string, caption string) error {
	return s.client.SendPhoto(chatID, photoURL, caption)
}

// SendDocument sends a document to Telegram
func (s *TelegramService) SendDocument(chatID int64, documentURL string, caption string) error {
	return s.client.SendDocument(chatID, documentURL, caption)
}

// SendLocation sends a location to Telegram
func (s *TelegramService) SendLocation(chatID int64, latitude float64, longitude float64) error {
	return s.client.SendLocation(chatID, latitude, longitude)
}

// ProcessWebhook processes incoming Telegram webhook
func (s *TelegramService) ProcessWebhook(updateData map[string]interface{}) error {
	// Convert map to Update struct
	update, err := s.parseUpdate(updateData)
	if err != nil {
		logger.Error("Failed to parse Telegram update", err)
		return err
	}

	// Save update to database
	if err := s.saveUpdate(update); err != nil {
		logger.Error("Failed to save Telegram update", err)
		return err
	}

	// Process the update
	return s.processUpdate(update)
}

func (s *TelegramService) parseUpdate(data map[string]interface{}) (*telegram.Update, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var update telegram.Update
	if err := json.Unmarshal(jsonData, &update); err != nil {
		return nil, err
	}

	return &update, nil
}

func (s *TelegramService) saveUpdate(update *telegram.Update) error {
	telegramMessage := &models.TelegramMessage{
		ID:        uuid.New(),
		UpdateID:  update.UpdateID,
		CreatedAt: time.Now(),
	}

	if update.Message != nil {
		telegramMessage.ChatID = update.Message.ChatID
		telegramMessage.Text = update.Message.Text
		telegramMessage.MessageType = "message"
	} else if update.CallbackQuery != nil {
		telegramMessage.ChatID = update.CallbackQuery.From.ID
		telegramMessage.Text = update.CallbackQuery.Data
		telegramMessage.MessageType = "callback_query"
	}

	return s.db.DB.Create(telegramMessage).Error
}

func (s *TelegramService) processUpdate(update *telegram.Update) error {
	// Handle auto-replies for Telegram
	if update.Message != nil {
		return s.handleTelegramMessage(update.Message)
	}

	// Handle callback queries
	if update.CallbackQuery != nil {
		return s.handleCallbackQuery(update.CallbackQuery)
	}

	return nil
}

func (s *TelegramService) handleTelegramMessage(message *telegram.Message) error {
	logger.Info("Handling Telegram message", map[string]interface{}{
		"chat_id": message.ChatID,
		"text":    message.Text,
	})

	// Get auto-replies for Telegram
	var autoReplies []models.AutoReply
	if err := s.db.DB.Where("is_active = ? AND platform = ?", true, "telegram").Find(&autoReplies).Error; err != nil {
		logger.Error("Failed to get Telegram auto-replies", err)
		return err
	}

	// Check for matching auto-reply
	for _, reply := range autoReplies {
		if s.matchesTrigger(message.Text, reply.Trigger, reply.MatchType) {
			// Send auto-reply
			if err := s.SendMessage(message.ChatID, reply.Response); err != nil {
				logger.Error("Failed to send auto-reply", err)
				return err
			}

			// Save auto-reply usage
			s.saveAutoReplyUsage(reply.ID, message.ChatID)
			return nil
		}
	}

	// Default response if no auto-reply matches
	return s.handleDefaultResponse(message)
}

func (s *TelegramService) handleCallbackQuery(callbackQuery *telegram.CallbackQuery) error {
	logger.Info("Handling Telegram callback query", map[string]interface{}{
		"callback_id": callbackQuery.ID,
		"data":        callbackQuery.Data,
	})

	// Answer the callback query
	responseText := fmt.Sprintf("You selected: %s", callbackQuery.Data)
	if err := s.client.AnswerCallbackQuery(callbackQuery.ID, responseText); err != nil {
		logger.Error("Failed to answer callback query", err)
		return err
	}

	// Handle specific callback data
	switch callbackQuery.Data {
	case "menu_main":
		return s.showMainMenu(callbackQuery.From.ID)
	case "menu_games":
		return s.showGamesMenu(callbackQuery.From.ID)
	case "menu_utilities":
		return s.showUtilitiesMenu(callbackQuery.From.ID)
	case "menu_business":
		return s.showBusinessMenu(callbackQuery.From.ID)
	default:
		// Handle game or utility callbacks
		if len(callbackQuery.Data) > 5 && callbackQuery.Data[:5] == "game_" {
			return s.handleGameCallback(callbackQuery)
		}
		if len(callbackQuery.Data) > 5 && callbackQuery.Data[:5] == "util_" {
			return s.handleUtilityCallback(callbackQuery)
		}
	}

	return nil
}

func (s *TelegramService) handleDefaultResponse(message *telegram.Message) error {
	// Show main menu with inline keyboard
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🎮 Games", "callback_data": "menu_games"},
				{"text": "🛠️ Utilities", "callback_data": "menu_utilities"},
			},
			{
				{"text": "💼 Business", "callback_data": "menu_business"},
				{"text": "📊 Analytics", "callback_data": "menu_analytics"},
			},
			{
				{"text": "⚙️ Settings", "callback_data": "menu_settings"},
				{"text": "❓ Help", "callback_data": "menu_help"},
			},
		},
	}

	welcomeText := "Welcome to Multi-Platform Bot! 🚀\n\nI can help you with:\n• Games and entertainment\n• Productivity tools\n• Business management\n• WhatsApp integration\n\nChoose an option below:"
	
	return s.SendMessageWithMarkup(message.ChatID, welcomeText, keyboard)
}

func (s *TelegramService) showMainMenu(chatID int64) error {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🎮 Games", "callback_data": "menu_games"},
				{"text": "🛠️ Utilities", "callback_data": "menu_utilities"},
			},
			{
				{"text": "💼 Business", "callback_data": "menu_business"},
				{"text": "📊 Analytics", "callback_data": "menu_analytics"},
			},
			{
				{"text": "⚙️ Settings", "callback_data": "menu_settings"},
				{"text": "❓ Help", "callback_data": "menu_help"},
			},
		},
	}

	text := "Main Menu 🏠\n\nWhat would you like to do?"
	return s.SendMessageWithMarkup(chatID, text, keyboard)
}

func (s *TelegramService) showGamesMenu(chatID int64) error {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🧠 Trivia", "callback_data": "game_trivia"},
				{"text": "🔢 Math Quiz", "callback_data": "game_math"},
			},
			{
				{"text": "📝 Word Game", "callback_data": "game_word"},
				{"text": "🧩 Memory Game", "callback_data": "game_memory"},
			},
			{
				{"text": "🏆 Leaderboard", "callback_data": "game_leaderboard"},
				{"text": "📊 Game Stats", "callback_data": "game_stats"},
			},
			{
				{"text": "🔙 Back to Main", "callback_data": "menu_main"},
			},
		},
	}

	text := "🎮 Games Menu\n\nChoose a game to play:"
	return s.SendMessageWithMarkup(chatID, text, keyboard)
}

func (s *TelegramService) showUtilitiesMenu(chatID int64) error {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "📱 QR Code", "callback_data": "util_qr"},
				{"text": "🔗 Short Link", "callback_data": "util_shortlink"},
			},
			{
				{"text": "💱 Currency", "callback_data": "util_currency"},
				{"text": "🌤️ Weather", "callback_data": "util_weather"},
			},
			{
				{"text": "🌐 Translate", "callback_data": "util_translate"},
				{"text": "📍 Location", "callback_data": "util_location"},
			},
			{
				{"text": "📊 Poll", "callback_data": "util_poll"},
				{"text": "⏰ Timer", "callback_data": "util_timer"},
			},
			{
				{"text": "🔙 Back to Main", "callback_data": "menu_main"},
			},
		},
	}

	text := "🛠️ Utilities Menu\n\nChoose a utility tool:"
	return s.SendMessageWithMarkup(chatID, text, keyboard)
}

func (s *TelegramService) showBusinessMenu(chatID int64) error {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "📦 Products", "callback_data": "business_products"},
				{"text": "📋 Orders", "callback_data": "business_orders"},
			},
			{
				{"text": "💰 Invoices", "callback_data": "business_invoices"},
				{"text": "👥 Customers", "callback_data": "business_customers"},
			},
			{
				{"text": "📊 Business Stats", "callback_data": "business_stats"},
				{"text": "📈 Sales Report", "callback_data": "business_report"},
			},
			{
				{"text": "🔙 Back to Main", "callback_data": "menu_main"},
			},
		},
	}

	text := "💼 Business Menu\n\nManage your business:"
	return s.SendMessageWithMarkup(chatID, text, keyboard)
}

func (s *TelegramService) handleGameCallback(callbackQuery *telegram.CallbackQuery) error {
	// Handle game-related callbacks
	gameType := callbackQuery.Data[5:] // Remove "game_" prefix
	
	switch gameType {
	case "trivia":
		return s.startTriviaGame(callbackQuery.From.ID)
	case "math":
		return s.startMathGame(callbackQuery.From.ID)
	case "word":
		return s.startWordGame(callbackQuery.From.ID)
	case "memory":
		return s.startMemoryGame(callbackQuery.From.ID)
	case "leaderboard":
		return s.showGameLeaderboard(callbackQuery.From.ID)
	case "stats":
		return s.showGameStats(callbackQuery.From.ID)
	default:
		return s.SendMessage(callbackQuery.From.ID, "Game not implemented yet!")
	}
}

func (s *TelegramService) handleUtilityCallback(callbackQuery *telegram.CallbackQuery) error {
	// Handle utility-related callbacks
	utilType := callbackQuery.Data[5:] // Remove "util_" prefix
	
	switch utilType {
	case "qr":
		return s.showQRInstructions(callbackQuery.From.ID)
	case "shortlink":
		return s.showShortLinkInstructions(callbackQuery.From.ID)
	case "currency":
		return s.showCurrencyInstructions(callbackQuery.From.ID)
	case "weather":
		return s.showWeatherInstructions(callbackQuery.From.ID)
	case "translate":
		return s.showTranslateInstructions(callbackQuery.From.ID)
	case "location":
		return s.showLocationInstructions(callbackQuery.From.ID)
	case "poll":
		return s.showPollInstructions(callbackQuery.From.ID)
	case "timer":
		return s.showTimerInstructions(callbackQuery.From.ID)
	default:
		return s.SendMessage(callbackQuery.From.ID, "Utility not implemented yet!")
	}
}

func (s *TelegramService) startTriviaGame(chatID int64) error {
	text := "🧠 Trivia Game\n\nStarting trivia game...\n\nPlease wait while I prepare questions for you!"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) startMathGame(chatID int64) error {
	text := "🔢 Math Quiz\n\nStarting math quiz...\n\nGet ready for some mathematical challenges!"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) startWordGame(chatID int64) error {
	text := "📝 Word Game\n\nStarting word game...\n\nLet's test your vocabulary!"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) startMemoryGame(chatID int64) error {
	text := "🧩 Memory Game\n\nStarting memory game...\n\nTest your memory skills!"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showGameLeaderboard(chatID int64) error {
	text := "🏆 Game Leaderboard\n\nTop players:\n1. Player1 - 1000 points\n2. Player2 - 850 points\n3. Player3 - 750 points"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showGameStats(chatID int64) error {
	text := "📊 Game Statistics\n\nYour stats:\n• Games played: 10\n• Total score: 5000\n• Average score: 500\n• Win rate: 70%"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showQRInstructions(chatID int64) error {
	text := "📱 QR Code Generator\n\nTo generate a QR code, send me a message like:\n`Generate QR for: https://example.com`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showShortLinkInstructions(chatID int64) error {
	text := "🔗 Short Link Generator\n\nTo create a short link, send me a message like:\n`Shorten: https://very-long-url-example.com`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showCurrencyInstructions(chatID int64) error {
	text := "💱 Currency Converter\n\nTo convert currency, send me a message like:\n`Convert 100 USD to EUR`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showWeatherInstructions(chatID int64) error {
	text := "🌤️ Weather Information\n\nTo get weather info, send me a message like:\n`Weather in Jakarta, Indonesia`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showTranslateInstructions(chatID int64) error {
	text := "🌐 Translation Service\n\nTo translate text, send me a message like:\n`Translate \"Hello\" to Indonesian`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showLocationInstructions(chatID int64) error {
	text := "📍 Location Services\n\nTo get location info, send me your location or:\n`Location info for: -6.2088, 106.8456`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showPollInstructions(chatID int64) error {
	text := "📊 Poll Creator\n\nTo create a poll, send me a message like:\n`Create poll: What's your favorite color? Options: Red, Blue, Green`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) showTimerInstructions(chatID int64) error {
	text := "⏰ Timer Service\n\nTo set a timer, send me a message like:\n`Timer for 5 minutes: Take a break`"
	return s.SendMessage(chatID, text)
}

func (s *TelegramService) matchesTrigger(text, trigger, matchType string) bool {
	switch matchType {
	case "exact":
		return text == trigger
	case "contains":
		return strings.Contains(text, trigger)
	case "regex":
		// Simple regex matching (in production, use proper regex library)
		return strings.Contains(text, trigger)
	default:
		return false
	}
}

func (s *TelegramService) saveAutoReplyUsage(replyID uuid.UUID, chatID int64) {
	usage := &models.AutoReplyUsage{
		ID:           uuid.New(),
		AutoReplyID:  replyID,
		Platform:     "telegram",
		RecipientID:  fmt.Sprintf("%d", chatID),
		UsedAt:       time.Now(),
	}

	if err := s.db.DB.Create(usage).Error; err != nil {
		logger.Error("Failed to save auto-reply usage", err)
	}
}

func (s *TelegramService) GetStatus() map[string]interface{} {
	return s.client.GetStatus()
}

func (s *TelegramService) StartPolling() error {
	return s.client.StartPolling(func(update telegram.Update) error {
		return s.processUpdate(&update)
	})
}

func (s *TelegramService) SetWebhook(webhookURL string) error {
	return s.client.SetWebhook(webhookURL)
}

func (s *TelegramService) DeleteWebhook() error {
	return s.client.DeleteWebhook()
}

func (s *TelegramService) GetWebhookInfo() (map[string]interface{}, error) {
	return s.client.GetWebhookInfo()
}

func (s *TelegramService) GetMe() (*telegram.User, error) {
	return s.client.GetMe()
}

func (s *TelegramService) GetUpdates(offset int) ([]telegram.Update, error) {
	return s.client.GetUpdates(offset)
}

func (s *TelegramService) AnswerCallbackQuery(callbackQueryID string, text string) error {
	return s.client.AnswerCallbackQuery(callbackQueryID, text)
}

func (s *TelegramService) SendPoll(chatID int64, question string, options []string) error {
	// Create inline keyboard for poll
	keyboard := map[string]interface{}{
		"inline_keyboard": make([][]map[string]interface{}, len(options)),
	}

	for i, option := range options {
		keyboard["inline_keyboard"].([][]map[string]interface{})[i] = []map[string]interface{}{
			{
				"text":          option,
				"callback_data": fmt.Sprintf("poll_vote_%d", i),
			},
		}
	}

	text := fmt.Sprintf("📊 Poll: %s\n\nClick an option to vote:", question)
	return s.SendMessageWithMarkup(chatID, text, keyboard)
}

func (s *TelegramService) SendBroadcast(recipients []int64, message string) error {
	successCount := 0
	failureCount := 0

	for _, chatID := range recipients {
		if err := s.SendMessage(chatID, message); err != nil {
			logger.Error("Failed to send broadcast message", err, map[string]interface{}{
				"chat_id": chatID,
			})
			failureCount++
		} else {
			successCount++
		}
	}

	logger.Info("Telegram broadcast completed", map[string]interface{}{
		"total":     len(recipients),
		"success":   successCount,
		"failure":   failureCount,
	})

	return nil
}

func (s *TelegramService) CreateTelegramUser(userID int64, username, firstName, lastName string) (*models.TelegramUser, error) {
	telegramUser := &models.TelegramUser{
		ID:         uuid.New(),
		TelegramID: userID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.DB.Create(telegramUser).Error; err != nil {
		logger.Error("Failed to create Telegram user", err)
		return nil, err
	}

	return telegramUser, nil
}

func (s *TelegramService) GetTelegramUser(userID int64) (*models.TelegramUser, error) {
	var user models.TelegramUser
	if err := s.db.DB.Where("telegram_id = ?", userID).First(&user).Error; err != nil {
		logger.Error("Failed to get Telegram user", err)
		return nil, err
	}

	return &user, nil
}

func (s *TelegramService) UpdateTelegramUser(userID int64, username, firstName, lastName string) (*models.TelegramUser, error) {
	var user models.TelegramUser
	if err := s.db.DB.Where("telegram_id = ?", userID).First(&user).Error; err != nil {
		logger.Error("Failed to get Telegram user for update", err)
		return nil, err
	}

	user.Username = username
	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&user).Error; err != nil {
		logger.Error("Failed to update Telegram user", err)
		return nil, err
	}

	return &user, nil
}

func (s *TelegramService) GetTelegramMessages(chatID int64, limit int) ([]models.TelegramMessage, error) {
	var messages []models.TelegramMessage
	if err := s.db.DB.Where("chat_id = ?", chatID).
		Order("created_at desc").
		Limit(limit).
		Find(&messages).Error; err != nil {
		logger.Error("Failed to get Telegram messages", err)
		return nil, err
	}

	return messages, nil
}

func (s *TelegramService) GetTelegramStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total messages
	var totalMessages int64
	s.db.DB.Model(&models.TelegramMessage{}).Count(&totalMessages)
	stats["total_messages"] = totalMessages

	// Total users
	var totalUsers int64
	s.db.DB.Model(&models.TelegramUser{}).Count(&totalUsers)
	stats["total_users"] = totalUsers

	// Messages today
	var todayMessages int64
	today := time.Now().Truncate(24 * time.Hour)
	s.db.DB.Model(&models.TelegramMessage{}).
		Where("created_at >= ?", today).
		Count(&todayMessages)
	stats["today_messages"] = todayMessages

	// Active users (last 30 days)
	var activeUsers int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	s.db.DB.Model(&models.TelegramMessage{}).
		Where("created_at >= ?", thirtyDaysAgo).
		Select("COUNT(DISTINCT chat_id)").
		Scan(&activeUsers)
	stats["active_users_30d"] = activeUsers

	return stats, nil
}

func (s *TelegramService) GetTelegramBroadcasts(userID uuid.UUID, status string, page int, limit int) ([]models.TelegramBroadcast, int, error) {
	// This would be implemented by TelegramBroadcastService
	// For now, return empty results
	return []models.TelegramBroadcast{}, 0, nil
}

func (s *TelegramService) CreateTelegramBroadcast(name string, message string, recipients []int64, userID uuid.UUID) (*models.TelegramBroadcast, error) {
	// This would be implemented by TelegramBroadcastService
	// For now, return a mock broadcast
	broadcast := &models.TelegramBroadcast{
		ID:         uuid.New(),
		Name:       name,
		Message:    message,
		Recipients: recipients,
		Status:     "created",
		UserID:     userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return broadcast, nil
}

func (s *TelegramService) SendTelegramBroadcast(broadcastID uuid.UUID) error {
	// This would be implemented by TelegramBroadcastService
	return nil
}

func (s *TelegramService) GetTelegramBroadcastStats(broadcastID uuid.UUID) (map[string]interface{}, error) {
	// This would be implemented by TelegramBroadcastService
	return map[string]interface{}{
		"broadcast_id": broadcastID,
		"status":        "completed",
		"total_sent":    100,
		"success_rate":  95,
	}, nil
}