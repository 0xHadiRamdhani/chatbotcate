package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Base model with common fields
type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" sql:"index"`
}

// User model
type User struct {
	BaseModel
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	PhoneNumber  string    `gorm:"unique"`
	DisplayName  string
	Avatar       string
	IsActive     bool      `gorm:"default:true"`
	IsAdmin      bool      `gorm:"default:false"`
	LastLoginAt  *time.Time
	Preferences  UserPreferences `gorm:"foreignKey:UserID"`
	Points       int           `gorm:"default:0"`
	Level        int           `gorm:"default:1"`
}

type UserPreferences struct {
	BaseModel
	UserID           uuid.UUID `gorm:"type:uuid;not null"`
	Language         string    `gorm:"default:'id'"`
	Timezone         string    `gorm:"default:'Asia/Jakarta'"`
	EnableGames      bool      `gorm:"default:true"`
	EnableBusiness   bool      `gorm:"default:true"`
	EnableUtils      bool      `gorm:"default:true"`
	EnableModeration bool      `gorm:"default:true"`
}

// Contact model
type Contact struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	PhoneNumber string    `gorm:"not null"`
	DisplayName string
	ProfilePic  string
	IsBlocked   bool `gorm:"default:false"`
	IsGroup     bool `gorm:"default:false"`
	GroupID     string
	LastMessage time.Time
}

// Message model
type Message struct {
	BaseModel
	UserID       uuid.UUID  `gorm:"type:uuid;not null"`
	ContactID    uuid.UUID  `gorm:"type:uuid;not null"`
	MessageID    string     `gorm:"unique;not null"`
	Content      string     `gorm:"type:text"`
	MessageType  string     `gorm:"not null"` // text, image, audio, video, document
	Direction    string     `gorm:"not null"` // incoming, outgoing
	Status       string     `gorm:"default:'sent'"` // sent, delivered, read, failed
	MediaURL     string
	MediaMimeType string
	Timestamp    time.Time
	IsForwarded  bool `gorm:"default:false"`
	IsReply      bool `gorm:"default:false"`
	ReplyToID    string
}

// AutoReply model
type AutoReply struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Keyword     string    `gorm:"not null"`
	Response    string    `gorm:"not null"`
	IsActive    bool      `gorm:"default:true"`
	MatchType   string    `gorm:"default:'exact'"` // exact, contains, regex
	ReplyType   string    `gorm:"default:'text'"`  // text, image, template
	MediaURL    string
	TemplateID  string
}

// Broadcast model
type Broadcast struct {
	BaseModel
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	Name        string     `gorm:"not null"`
	Content     string     `gorm:"not null"`
	MessageType string     `gorm:"default:'text'"`
	MediaURL    string
	Recipients  []BroadcastRecipient
	Status      string    `gorm:"default:'draft'"` // draft, scheduled, sending, sent, failed
	ScheduledAt *time.Time
	SentAt      *time.Time
	TotalSent   int       `gorm:"default:0"`
	TotalFailed int       `gorm:"default:0"`
}

type BroadcastRecipient struct {
	BaseModel
	BroadcastID uuid.UUID `gorm:"type:uuid;not null"`
	ContactID   uuid.UUID `gorm:"type:uuid;not null"`
	Status      string    `gorm:"default:'pending'"` // pending, sent, failed
	SentAt      *time.Time
	Error       string
}

// Group model
type Group struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	GroupID     string    `gorm:"unique;not null"`
	Name        string    `gorm:"not null"`
	Description string
	ProfilePic  string
	Members     []GroupMember
	IsActive    bool `gorm:"default:true"`
}

type GroupMember struct {
	BaseModel
	GroupID   uuid.UUID `gorm:"type:uuid;not null"`
	ContactID uuid.UUID `gorm:"type:uuid;not null"`
	Role      string    `gorm:"default:'member'"` // admin, member
	JoinedAt  time.Time
}

// Game models
type GameScore struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	GameType  string    `gorm:"not null"` // quiz, tebak_gambar, math_challenge
	Score     int       `gorm:"default:0"`
	Level     int       `gorm:"default:1"`
	PlayCount int       `gorm:"default:0"`
	BestScore int       `gorm:"default:0"`
}

type Quiz struct {
	BaseModel
	Question    string `gorm:"not null"`
	Options     string `gorm:"not null"` // JSON array
	CorrectAnswer int    `gorm:"not null"`
	Category    string
	Difficulty  string `gorm:"default:'easy'"` // easy, medium, hard
	Points      int    `gorm:"default:10"`
}

type QuizSession struct {
	BaseModel
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	QuizID     uuid.UUID `gorm:"type:uuid;not null"`
	Score      int       `gorm:"default:0"`
	CurrentQuestion int  `gorm:"default:0"`
	TotalQuestions int   `gorm:"default:0"`
	Status     string    `gorm:"default:'active'"` // active, completed, abandoned
	StartedAt  time.Time
	CompletedAt *time.Time
}

// Business models
type Product struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	Currency    string    `gorm:"default:'IDR'"`
	Category    string
	Images      string    `gorm:"type:text"` // JSON array
	Stock       int       `gorm:"default:0"`
	IsActive    bool      `gorm:"default:true"`
}

type Order struct {
	BaseModel
	UserID       uuid.UUID      `gorm:"type:uuid;not null"`
	ContactID    uuid.UUID      `gorm:"type:uuid;not null"`
	OrderNumber  string         `gorm:"unique;not null"`
	Status       string         `gorm:"default:'pending'"` // pending, confirmed, processing, shipped, delivered, cancelled
	TotalAmount  float64        `gorm:"not null"`
	Currency     string         `gorm:"default:'IDR'"`
	Items        []OrderItem
	ShippingAddress string      `gorm:"type:text"`
	Notes        string         `gorm:"type:text"`
	PaidAt       *time.Time
	ShippedAt    *time.Time
	DeliveredAt  *time.Time
}

type OrderItem struct {
	BaseModel
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Subtotal  float64   `gorm:"not null"`
}

// Reminder model
type Reminder struct {
	BaseModel
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	ContactID   uuid.UUID  `gorm:"type:uuid;not null"`
	Title       string     `gorm:"not null"`
	Description string     `gorm:"type:text"`
	RemindAt    time.Time  `gorm:"not null"`
	IsRecurring bool       `gorm:"default:false"`
	RecurringType string   `gorm:"default:'none'"` // none, daily, weekly, monthly
	Status      string     `gorm:"default:'active'"` // active, completed, cancelled
	CompletedAt *time.Time
}

// Moderation models
type BlockedWord struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Word   string    `gorm:"not null"`
	Action string    `gorm:"default:'block'"` // block, replace, warn
	ReplaceWith string
}

type SpamReport struct {
	BaseModel
	ReporterID uuid.UUID `gorm:"type:uuid;not null"`
	ReportedID uuid.UUID `gorm:"type:uuid;not null"`
	Reason     string    `gorm:"not null"`
	Evidence   string    `gorm:"type:text"`
	Status     string    `gorm:"default:'pending'"` // pending, reviewed, resolved
}

// Utility models
type WeatherData struct {
	BaseModel
	City      string    `gorm:"not null"`
	Country   string    `gorm:"not null"`
	Temp      float64   `gorm:"not null"`
	FeelsLike float64   `gorm:"not null"`
	Humidity  int       `gorm:"not null"`
	Pressure  int       `gorm:"not null"`
	WindSpeed float64   `gorm:"not null"`
	Condition string    `gorm:"not null"`
	FetchedAt time.Time `gorm:"not null"`
}

type CurrencyRate struct {
	BaseModel
	FromCurrency string    `gorm:"not null"`
	ToCurrency   string    `gorm:"not null"`
	Rate         float64   `gorm:"not null"`
	FetchedAt    time.Time `gorm:"not null"`
}

// Analytics model
type Analytics struct {
	BaseModel
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	MetricType     string    `gorm:"not null"` // messages_sent, games_played, orders_created
	MetricValue    int       `gorm:"default:0"`
	AdditionalData string    `gorm:"type:text"` // JSON for extra data
}

// Custom command model
type CustomCommand struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Command     string    `gorm:"not null"`
	Description string
	Response    string    `gorm:"not null"`
	IsActive    bool      `gorm:"default:true"`
	TriggerType string    `gorm:"default:'exact'"` // exact, contains, regex
}

// Template model
type Template struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"not null"`
	Content     string    `gorm:"not null"`
	Variables   string    `gorm:"type:text"` // JSON array of variable names
	Category    string
	IsActive    bool `gorm:"default:true"`
}

// System log model
type SystemLog struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid"`
	Level       string    `gorm:"not null"` // info, warn, error, debug
	Message     string    `gorm:"not null;type:text"`
	Context     string    `gorm:"type:text"` // JSON for additional context
	IPAddress   string
	UserAgent   string
}