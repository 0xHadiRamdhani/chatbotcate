package models

import (
	"time"

	"github.com/google/uuid"
)

// TelegramUser represents a Telegram user
type TelegramUser struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	TelegramID int64      `json:"telegram_id" gorm:"uniqueIndex"`
	Username   string     `json:"username"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TelegramMessage represents a message sent/received via Telegram
type TelegramMessage struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UpdateID    int        `json:"update_id"`
	ChatID      int64      `json:"chat_id"`
	MessageID   int        `json:"message_id"`
	Text        string     `json:"text" gorm:"type:text"`
	MessageType string     `json:"message_type"` // message, callback_query, inline_query
	FromUserID  int64      `json:"from_user_id"`
	FromUsername string    `json:"from_username"`
	Direction   string     `json:"direction"` // incoming, outgoing
	Status      string     `json:"status"`    // sent, delivered, read
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TelegramBroadcast represents a broadcast message for Telegram
type TelegramBroadcast struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primary_key"`
	Name       string      `json:"name"`
	Message    string      `json:"message" gorm:"type:text"`
	Recipients []int64     `json:"recipients" gorm:"type:jsonb"`
	Status     string      `json:"status"` // pending, sending, completed, failed
	SuccessCount int       `json:"success_count"`
	FailureCount int       `json:"failure_count"`
	UserID     uuid.UUID   `json:"user_id" gorm:"type:uuid"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	SentAt     *time.Time  `json:"sent_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// TelegramAutoReply represents auto-reply rules for Telegram
type TelegramAutoReply struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name      string     `json:"name"`
	Trigger   string     `json:"trigger"`
	Response  string     `json:"response" gorm:"type:text"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	MatchType string     `json:"match_type"` // exact, contains, regex
	Keywords  []string   `json:"keywords" gorm:"type:jsonb"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TelegramGroup represents a Telegram group/chat
type TelegramGroup struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	ChatID      int64      `json:"chat_id" gorm:"uniqueIndex"`
	Title       string     `json:"title"`
	Type        string     `json:"type"` // private, group, supergroup, channel
	Description string     `json:"description"`
	MemberCount int        `json:"member_count"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TelegramBot represents bot configuration for Telegram
type TelegramBot struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name         string     `json:"name"`
	Username     string     `json:"username"`
	APIKey       string     `json:"api_key"`
	WebhookURL   string     `json:"webhook_url"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastError    string     `json:"last_error"`
	LastErrorAt  *time.Time `json:"last_error_at"`
	MessageCount int64      `json:"message_count" gorm:"default:0"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TelegramCommand represents custom commands for Telegram bot
type TelegramCommand struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Command     string     `json:"command" gorm:"uniqueIndex"`
	Description string     `json:"description"`
	Response    string     `json:"response" gorm:"type:text"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TelegramInlineKeyboard represents inline keyboard markup
type TelegramInlineKeyboard struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name      string     `json:"name"`
	Buttons   []Button   `json:"buttons" gorm:"type:jsonb"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Button represents a button in inline keyboard
type Button struct {
	Text                         string `json:"text"`
	CallbackData                 string `json:"callback_data,omitempty"`
	URL                          string `json:"url,omitempty"`
	SwitchInlineQuery            string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
}

// TelegramWebhook represents webhook configuration
type TelegramWebhook struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	URL          string     `json:"url"`
	SecretToken  string     `json:"secret_token"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastUpdateID int        `json:"last_update_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TelegramAnalytics represents analytics data for Telegram
type TelegramAnalytics struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Date            time.Time  `json:"date"`
	MessagesSent    int64      `json:"messages_sent"`
	MessagesReceived int64     `json:"messages_received"`
	UsersInteracted int64      `json:"users_interacted"`
	GroupsInteracted int64     `json:"groups_interacted"`
	CommandsUsed    int64      `json:"commands_used"`
	ErrorsCount     int64      `json:"errors_count"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TelegramSession represents user session data
type TelegramSession struct {
	ID           uuid.UUID              `json:"id" gorm:"type:uuid;primary_key"`
	UserID       int64                  `json:"user_id"`
	ChatID       int64                  `json:"chat_id"`
	SessionData  map[string]interface{} `json:"session_data" gorm:"type:jsonb"`
	LastActivity time.Time             `json:"last_activity"`
	IsActive     bool                   `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// TelegramNotification represents notification settings
type TelegramNotification struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID       int64      `json:"user_id"`
	ChatID       int64      `json:"chat_id"`
	Type         string     `json:"type"` // message, broadcast, reminder, etc.
	IsEnabled    bool       `json:"is_enabled" gorm:"default:true"`
	QuietHours   string     `json:"quiet_hours"` // e.g., "22:00-08:00"
	Language     string     `json:"language" gorm:"default:'en'"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TelegramIntegration represents integration settings between platforms
type TelegramIntegration struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID           uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	TelegramUserID   int64      `json:"telegram_user_id"`
	WhatsAppUserID   string     `json:"whatsapp_user_id"`
	IsSyncEnabled    bool       `json:"is_sync_enabled" gorm:"default:true"`
	SyncDirection    string     `json:"sync_direction"` // telegram_to_whatsapp, whatsapp_to_telegram, both
	LastSyncAt       *time.Time `json:"last_sync_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}